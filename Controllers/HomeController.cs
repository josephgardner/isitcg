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
        public IActionResult Submit(string ingredients)
        {
            var results = _ingredientHandler.CreateResults(ingredients);
            var id = _fileManager.Write(results);
            return RedirectToAction("Results", new { id });
        }

        public IActionResult Results(string id)
        {
            var results = _fileManager.Read(id);
            ViewData["result"] = results.Matches.Any() ? "NO!" : "YES!";
            return View(results);
        }

        public IActionResult Error()
        {
            return View();
        }
    }
}
