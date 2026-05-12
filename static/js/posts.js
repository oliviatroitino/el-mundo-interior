class PostsManager {
  constructor() {
    this.myList    = document.getElementById('my-posts')
    this.otherList = document.getElementById('other-posts')
    this.template  = document.getElementById('post-tpl')
    this.form      = document.querySelector('.post-form__form')
    this.worldSlug   = this.myList?.dataset.world
    this.sectionSlug = this.myList?.dataset.section ?? ''
  }

  init() {
    if (!this.myList) return
    this.errorBanner = document.getElementById('post-error')
    this.loadPosts()
    this.form?.addEventListener('submit', (e) => this.handleSubmit(e))

    // Botón de archivo: abre el input file oculto
    const fileInput  = this.form?.querySelector('#post-media')
    this.form?.querySelector('.post-form__upload')
      ?.addEventListener('click', () => fileInput?.click())

    // Botón de ubicación: muestra/oculta el input aplicando una clase CSS
    const locInput = this.form?.querySelector('.post-form__location-input')
    this.form?.querySelector('.post-form__location')
      ?.addEventListener('click', () => locInput?.classList.toggle('is-visible'))
  }

  // Muestra un mensaje de error al usuario durante 4 segundos.
  showError(message) {
    this.errorBanner.textContent = message
    this.errorBanner.hidden = false
    clearTimeout(this._errorTimer)
    this._errorTimer = setTimeout(() => { this.errorBanner.hidden = true }, 4000)
  }

  // Pide los posts a la API y los renderiza en sus contenedores.
  // El elemento "cargando" aparece antes del fetch y se elimina en finally,
  // tanto si hay éxito como si hay error.
  async loadPosts() {
    const loading = document.createElement('li')
    loading.textContent = 'Cargando...'
    this.myList.appendChild(loading)

    try {
      const url   = this.sectionSlug
        ? `/api/posts?world=${this.worldSlug}&section=${this.sectionSlug}`
        : `/api/posts?world=${this.worldSlug}`
      const res   = await fetch(url)
      const posts = await res.json()

      posts.forEach(post => {
        const card = this.createCard(post)
        post.mine ? this.myList.appendChild(card) : this.otherList.appendChild(card)
      })
    } catch (err) {
      this.showError('Error al cargar las publicaciones. Inténtalo de nuevo.')
    } finally {
      loading.remove()
      this.showEmptyState(this.myList, 'Aún no has publicado nada en este mundo.')
      this.showEmptyState(this.otherList, 'Aún no hay publicaciones de otros usuarios.')
    }
  }

  // Muestra un mensaje en la lista si está vacía tras cargar los posts.
  showEmptyState(list, message) {
    if (list.children.length === 0) {
      const empty = document.createElement('li')
      empty.className = 'post-list__empty'
      empty.textContent = message
      list.appendChild(empty)
    }
  }

  // Clona el <template>, rellena los huecos con los datos del post
  // y registra los botones de editar/borrar si el post es del usuario.
  createCard(post) {
    const clone = this.template.content.cloneNode(true)
    const li    = clone.querySelector('li')

    li.dataset.id      = post.id
    li.dataset.section = post.section_slug ?? ''
    li.dataset.location = post.location ?? ''

    clone.querySelector('.post__user').textContent = '@' + post.user_name
    clone.querySelector('.post__text').textContent = post.body
    clone.querySelector('.post__date').textContent = post.date

    if (post.location) {
      const loc       = clone.querySelector('.post__location')
      loc.textContent = post.location
      loc.hidden      = false
    }

    if (post.media_path) {
      const img  = clone.querySelector('.post__media')
      img.src    = post.media_path
      img.hidden = false
    }

    if (post.mine) {
      clone.querySelector('.post__actions').hidden = false
      clone.querySelector('.btn--edit').addEventListener('click', () => this.handleEdit(li))
      clone.querySelector('.btn--delete').addEventListener('click', () => this.handleDelete(li))
    }

    return clone
  }

  // Intercepta el envío del formulario y crea el post vía API
  // sin recargar la página. El botón queda deshabilitado mientras espera.
  async handleSubmit(e) {
    e.preventDefault()
    const form      = e.target
    const body      = form.querySelector('[name="body"]').value.trim()
    const section   = form.querySelector('[name="section_slug"]').value
    const location  = form.querySelector('[name="location"]').value.trim()
    const submitBtn = form.querySelector('.post-form__submit')

    if (!body) return

    submitBtn.disabled = true
    try {
      const res = await fetch('/api/posts', {
        method:  'POST',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify({
          world_slug:   this.worldSlug,
          section_slug: section,
          body,
          location,
        }),
      })
      if (!res.ok) throw new Error('Error del servidor')

      const post = await res.json()
      this.myList.prepend(this.createCard(post))
      form.reset()
    } catch (err) {
      this.showError('Error al publicar. Inténtalo de nuevo.')
    } finally {
      submitBtn.disabled = false
    }
  }

  // Clona el <template> de edición, rellena los campos con los valores actuales
  // y registra el botón guardar para llamar a PATCH y restaurar la card.
  handleEdit(li) {
    const id           = li.dataset.id
    const article      = li.querySelector('article')
    const originalHTML = article.innerHTML

    const clone      = document.getElementById('post-edit-tpl').content.cloneNode(true)
    const textarea   = clone.querySelector('.post-form__text')
    const locInput   = clone.querySelector('.post-form__location-input')
    const sectionSel = clone.querySelector('.post-form__section-select')
    const saveBtn    = clone.querySelector('.btn--save')

    textarea.value = li.querySelector('.post__text')?.textContent ?? ''
    locInput.value = li.dataset.location

    const sourceSelect = this.form?.querySelector('[name="section_slug"]')
                      ?? document.getElementById('sections-data')
    if (sourceSelect) {
      Array.from(sourceSelect.options).forEach(opt => {
        const o       = document.createElement('option')
        o.value       = opt.value
        o.textContent = opt.textContent
        if (opt.value === li.dataset.section) o.selected = true
        sectionSel.appendChild(o)
      })
    } else {
      sectionSel.remove()
    }

    article.innerHTML = ''
    article.appendChild(clone)
    textarea.focus()

    saveBtn.addEventListener('click', async () => {
      const newBody     = textarea.value.trim()
      const newLocation = locInput.value.trim()
      const newSection  = sectionSel.isConnected ? sectionSel.value : li.dataset.section

      if (!newBody) return

      saveBtn.disabled = true
      try {
        const res = await fetch(`/api/posts/${id}`, {
          method:  'PATCH',
          headers: { 'Content-Type': 'application/json' },
          body:    JSON.stringify({ body: newBody, location: newLocation, section_slug: newSection }),
        })
        if (!res.ok) throw new Error('Error del servidor')
        const updated = await res.json()

        li.dataset.section  = updated.section_slug ?? ''
        li.dataset.location = updated.location ?? ''

        article.innerHTML = originalHTML
        article.querySelector('.post__text').textContent = updated.body
        const locEl       = article.querySelector('.post__location')
        locEl.textContent = updated.location ?? ''
        locEl.hidden      = !updated.location

        article.querySelector('.btn--edit').addEventListener('click', () => this.handleEdit(li))
        article.querySelector('.btn--delete').addEventListener('click', () => this.handleDelete(li))
      } catch (err) {
        this.showError('Error al editar la publicación. Inténtalo de nuevo.')
      } finally {
        saveBtn.disabled = false
      }
    })
  }

  // Llama a DELETE y elimina la card del DOM si el servidor responde ok.
  async handleDelete(li) {
    const deleteBtn    = li.querySelector('.btn--delete')
    deleteBtn.disabled = true

    try {
      const res = await fetch(`/api/posts/${li.dataset.id}`, { method: 'DELETE' })
      if (!res.ok) throw new Error('Error del servidor')
      li.remove()
    } catch (err) {
      this.showError('Error al borrar la publicación. Inténtalo de nuevo.')
      deleteBtn.disabled = false
    }
  }
}

const manager = new PostsManager()
manager.init()
