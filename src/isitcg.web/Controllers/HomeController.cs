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
        private readonly IDatabase _database;
        public HomeController(IIngredientHandler ingredientHandler, IDatabase database)
        {
            if (ingredientHandler == null)
                throw new ArgumentNullException(nameof(ingredientHandler));
            if (database == null)
                throw new ArgumentNullException(nameof(database));

            _ingredientHandler = ingredientHandler;
            _database = database;
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
            _database.SortedSetIncrement(
                "products", results.ProductName, 1, CommandFlags.FireAndForget);
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
            var json = (await _database.HashGetAllAsync(id)).Last().Value;
            var results = JsonConvert.DeserializeObject<MatchResults>(json);
            return View(results);
        }

        public IActionResult Error()
        {
            return View();
        }
    }
}
