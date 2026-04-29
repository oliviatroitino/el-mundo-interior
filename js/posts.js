class PostsManager {
  constructor() {
    this.myList    = document.getElementById('my-posts')
    this.otherList = document.getElementById('other-posts')
    this.template  = document.getElementById('post-tpl')
    this.form      = document.querySelector('.post-form__form')
    this.worldSlug = this.myList?.dataset.world
  }

  init() {
    if (!this.myList) return
    this.loadPosts()
    this.form?.addEventListener('submit', (e) => this.handleSubmit(e))
  }

  // Pide los posts a la API y los renderiza en sus contenedores.
  // El elemento "cargando" aparece antes del fetch y se elimina en finally,
  // tanto si hay éxito como si hay error.
  async loadPosts() {
    const loading = document.createElement('li')
    loading.textContent = 'Cargando...'
    this.myList.appendChild(loading)

    try {
      const res   = await fetch(`/api/posts?world=${this.worldSlug}`)
      const posts = await res.json()

      posts.forEach(post => {
        const card = this.createCard(post)
        post.mine ? this.myList.appendChild(card) : this.otherList.appendChild(card)
      })
    } catch (err) {
      console.error('Error cargando posts:', err)
    } finally {
      loading.remove()
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
      console.error('Error creando post:', err)
    } finally {
      submitBtn.disabled = false
    }
  }

  // Reemplaza el contenido de la card por un formulario igual al de crear posts,
  // con los valores actuales pre-rellenados. Al guardar llama a PATCH y restaura la card.
  handleEdit(li) {
    const id          = li.dataset.id
    const article     = li.querySelector('article')
    const originalHTML = article.innerHTML

    // — Textarea con el texto actual —
    const textarea     = document.createElement('textarea')
    textarea.className = 'post-form__text'
    textarea.value     = li.querySelector('.post__text')?.textContent ?? ''

    // — Input de ubicación siempre visible (en edición no tiene sentido ocultarlo) —
    const locInput       = document.createElement('input')
    locInput.type        = 'text'
    locInput.className   = 'post-form__location-input'
    locInput.placeholder = 'Ubicación...'
    locInput.value       = li.dataset.location
    locInput.style.display = 'block'

    // — Select de sección (copia opciones del formulario principal) —
    const sectionSelect     = document.createElement('select')
    sectionSelect.className = 'post-form__section-select'
    const sourceSelect      = this.form?.querySelector('[name="section_slug"]')
    if (sourceSelect) {
      Array.from(sourceSelect.options).forEach(opt => {
        const o       = document.createElement('option')
        o.value       = opt.value
        o.textContent = opt.textContent
        if (opt.value === li.dataset.section) o.selected = true
        sectionSelect.appendChild(o)
      })
    }

    // — Botón guardar con texto explícito —
    const saveBtn         = document.createElement('button')
    saveBtn.type          = 'button'
    saveBtn.className     = 'btn btn--save'
    saveBtn.textContent   = 'Guardar'

    const options     = document.createElement('div')
    options.className = 'post-form__options'
    options.append(locInput, sectionSelect, saveBtn)

    const wrapper     = document.createElement('div')
    wrapper.className = 'post-form__input'
    wrapper.append(textarea, options)

    article.innerHTML = ''
    article.appendChild(wrapper)
    textarea.focus()

    saveBtn.addEventListener('click', async () => {
      const newBody     = textarea.value.trim()
      const newLocation = locInput.value.trim()
      const newSection  = sectionSelect.value

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

        // Restaurar la card con los nuevos valores
        article.innerHTML = originalHTML
        article.querySelector('.post__text').textContent = updated.body
        const locEl    = article.querySelector('.post__location')
        locEl.textContent = updated.location ?? ''
        locEl.hidden      = !updated.location

        // Reconectar listeners (innerHTML los elimina)
        article.querySelector('.btn--edit').addEventListener('click', () => this.handleEdit(li))
        article.querySelector('.btn--delete').addEventListener('click', () => this.handleDelete(li))
      } catch (err) {
        console.error('Error actualizando post:', err)
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
      console.error('Error borrando post:', err)
      deleteBtn.disabled = false
    }
  }
}

const manager = new PostsManager()
manager.init()
