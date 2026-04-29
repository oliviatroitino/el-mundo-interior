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

  // Sustituye el contenido del post por un formulario de edición inline
  // con textarea (body), campo de ubicación y selector de sección.
  // Al guardar llama a PATCH y restaura la vista con los nuevos valores.
  handleEdit(li) {
    const id      = li.dataset.id
    const textEl  = li.querySelector('.post__text')
    const locEl   = li.querySelector('.post__location')
    const editBtn = li.querySelector('.btn--edit')

    // — Textarea para el cuerpo —
    const textarea     = document.createElement('textarea')
    textarea.className = 'post__edit-input'
    textarea.value     = textEl.textContent
    textEl.replaceWith(textarea)

    // — Input para la ubicación —
    const locInput     = document.createElement('input')
    locInput.type      = 'text'
    locInput.className = 'post__edit-location'
    locInput.placeholder = 'Ubicación...'
    locInput.value     = li.dataset.location
    locEl.replaceWith(locInput)
    locInput.hidden = false

    // — Select para la sección (clona opciones del formulario principal) —
    const sourceSelect  = this.form?.querySelector('[name="section_slug"]')
    const sectionSelect = document.createElement('select')
    sectionSelect.className = 'post__edit-section'
    if (sourceSelect) {
      Array.from(sourceSelect.options).forEach(opt => {
        const o       = document.createElement('option')
        o.value       = opt.value
        o.textContent = opt.textContent
        if (opt.value === li.dataset.section) o.selected = true
        sectionSelect.appendChild(o)
      })
    }
    textarea.after(sectionSelect)

    // — Botón Guardar (reemplaza Editar) —
    const saveBtn       = document.createElement('button')
    saveBtn.type        = 'button'
    saveBtn.className   = 'btn btn--save'
    saveBtn.textContent = 'Guardar'
    editBtn.replaceWith(saveBtn)

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
          body:    JSON.stringify({
            body:         newBody,
            location:     newLocation,
            section_slug: newSection,
          }),
        })
        if (!res.ok) throw new Error('Error del servidor')

        const updated = await res.json()

        // Restaurar el texto actualizado
        const newText       = document.createElement('p')
        newText.className   = 'post__text'
        newText.textContent = updated.body
        textarea.replaceWith(newText)

        // Restaurar la ubicación actualizada
        const newLoc       = document.createElement('p')
        newLoc.className   = 'post__location'
        newLoc.textContent = updated.location ?? ''
        newLoc.hidden      = !updated.location
        locInput.replaceWith(newLoc)

        // Guardar nuevos valores en data attributes
        li.dataset.section  = updated.section_slug ?? ''
        li.dataset.location = updated.location ?? ''

        sectionSelect.remove()

        // Restaurar botón Editar
        const newEditBtn       = document.createElement('button')
        newEditBtn.type        = 'button'
        newEditBtn.className   = 'btn btn--edit'
        newEditBtn.textContent = 'Editar'
        newEditBtn.addEventListener('click', () => this.handleEdit(li))
        saveBtn.replaceWith(newEditBtn)
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
