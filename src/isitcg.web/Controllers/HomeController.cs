using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using StackExchange.Redis;
using Newtonsoft.Json;

namespace isitcg.Controllers
{
    public class HomeController : Controller
    {
        private readonly IIngredientHandler _ingredientHandler;
        private readonly ResultCounter _resultCounter;
        private readonly IDatabase _db;

        public HomeController(
            IIngredientHandler ingredientHandler,
            ResultCounter resultCounter,
            IDatabase db)
        {
            _ingredientHandler = ingredientHandler ?? throw new ArgumentNullException(nameof(ingredientHandler));
            _resultCounter = resultCounter ?? throw new ArgumentNullException(nameof(resultCounter));
            _db = db ?? throw new ArgumentNullException(nameof(db));
        }
        public IActionResult Index(Product product)
        {
            return View(product ?? new Product());
        }

        [HttpPost]
        [ValidateAntiForgeryToken]
        public IActionResult Submit(string productname, string ingredients)
        {
            var hash = _ingredientHandler.CreateHash(productname, ingredients);
            return RedirectToAction("ViewHash", new { hash });
        }

        [Route("view/{hash}")]
        public IActionResult ViewHash(string hash)
        {
            var results = _ingredientHandler.ResultsFromHash(hash);
            _resultCounter.Count(results);

            return View("Results", results);
        }

        [Route("edit/{hash}")]
        public IActionResult EditHash(string hash)
        {
            var product = _ingredientHandler.ProductFromHash(hash);
            return View("Index", product);
        }

        // Legacy results handler; read from redis. 
        public async Task<IActionResult> Results(string id)
        {
            var json = (await _db.HashGetAllAsync(id)).Last().Value;
            var results = JsonConvert.DeserializeObject<MatchResults>(json);
            return View(results);
        }

        public IActionResult Error()
        {
            return View();
        }
    }
}
