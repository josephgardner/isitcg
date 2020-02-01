using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;

namespace isitcg.Controllers
{
    public class HomeController : Controller
    {
        private readonly IIngredientHandler _ingredientHandler;
        private readonly IFileManager _fileManager;
        public HomeController(IIngredientHandler ingredientHandler, IFileManager fileManager)
        {
            if (ingredientHandler == null)
                throw new ArgumentNullException(nameof(ingredientHandler));
            if (fileManager == null)
                throw new ArgumentNullException(nameof(fileManager));

            _ingredientHandler = ingredientHandler;
            _fileManager = fileManager;
        }
        public IActionResult Index()
        {
            return View();
        }

        [HttpPost]
        [ValidateAntiForgeryToken]
        public IActionResult Submit(string productname, string ingredients)
        {
            var hash = _ingredientHandler.CreateHash(productname, ingredients);
            return RedirectToAction("Hash", new { hash });
        }

        public IActionResult Hash(string hash)
        {
            var results = _ingredientHandler.ResultsFromHash(hash);
            return View("Results", results);
        }

        // Legacy results handler; read from redis. 
        public async Task<IActionResult> Results(string id)
        {
            var results = await _fileManager.ReadAsync(id);
            return View(results);
        }

        public IActionResult Error()
        {
            return View();
        }
    }
}
