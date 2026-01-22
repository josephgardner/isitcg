// Simple mobile menu toggle - replaces dlmenu.js + malina.js
(function() {
  document.addEventListener('DOMContentLoaded', function() {
    var trigger = document.querySelector('.dl-trigger');
    var menu = document.querySelector('.dl-menu');

    if (trigger && menu) {
      trigger.addEventListener('click', function(e) {
        e.preventDefault();
        menu.classList.toggle('dl-menuopen');
        trigger.classList.toggle('dl-active');
      });

      // Close menu when clicking outside
      document.addEventListener('click', function(e) {
        if (!e.target.closest('.dl-menuwrapper') && menu.classList.contains('dl-menuopen')) {
          menu.classList.remove('dl-menuopen');
          trigger.classList.remove('dl-active');
        }
      });
    }
  });
})();
