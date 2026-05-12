(function () {
  const btn  = document.querySelector('.question-generator__more-link')
  const list = document.querySelector('.question-generator__list')
  const tpl  = document.getElementById('question-tpl')

  if (!btn || !list || !tpl) return

  btn.addEventListener('click', async () => {
    btn.disabled = true
    try {
      const res            = await fetch('/api/questions')
      const { question }   = await res.json()

      const clone = tpl.content.cloneNode(true)
      clone.querySelector('.question-generator__question').textContent = question
      list.insertBefore(clone, list.querySelector('.question-generator__more'))
    } catch (err) {
      console.error('Error al cargar pregunta:', err)
    } finally {
      btn.disabled = false
    }
  })
})()
