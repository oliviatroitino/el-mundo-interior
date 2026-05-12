(function () {
  const btn  = document.querySelector('.question-generator__more-link')
  const list = document.querySelector('.question-generator__list')
  const tpl  = document.getElementById('question-tpl')

  if (!btn || !list || !tpl) return

  btn.addEventListener('click', async () => {
    btn.disabled = true
    try {
      const res             = await fetch('/api/questions?count=4')
      const { questions }   = await res.json()

      list.querySelectorAll('.question-generator__item').forEach(el => el.remove())

      const more = list.querySelector('.question-generator__more')
      questions.forEach(question => {
        const clone = tpl.content.cloneNode(true)
        clone.querySelector('.question-generator__question').textContent = question
        list.insertBefore(clone, more)
      })
    } catch (err) {
      console.error('Error al cargar preguntas:', err)
    } finally {
      btn.disabled = false
    }
  })
})()
