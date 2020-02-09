using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Rewrite;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using StackExchange.Redis;

namespace isitcg
{
    public class Startup
    {
        public Startup(IHostingEnvironment env)
        {
            var builder = new ConfigurationBuilder()
                .SetBasePath(env.ContentRootPath)
                .AddJsonFile("appsettings.json", optional: false, reloadOnChange: true)
                .AddJsonFile($"appsettings.{env.EnvironmentName}.json", optional: true)
                .AddJsonFile("ingredientrules.json", optional: false, reloadOnChange: true)
                .AddEnvironmentVariables();
            Configuration = builder.Build();
        }

        public IConfigurationRoot Configuration { get; }

        // This method gets called by the runtime. Use this method to add services to the container.
        public void ConfigureServices(IServiceCollection services)
        {
            services.AddOptions();
            services.Configure<IngredientRules>(Configuration);
            services.AddTransient<IIngredientHandler, DefaultIngredientHandler>();

            var redisUri = new Uri(Configuration
                            .GetSection("REDISCLOUD_URL").Value);
            var addresses = Dns.GetHostAddressesAsync(redisUri.Host).Result;
            var ip = addresses[0].MapToIPv4().ToString();
            var password = redisUri.UserInfo.Split(':')[1];
            var connect = $"{ip}:{redisUri.Port},password={password}";
            var redis = ConnectionMultiplexer.Connect(connect);
            services.AddSingleton<IConnectionMultiplexer>(redis);
            services.AddTransient<IDatabase>(c =>
                c.GetRequiredService<IConnectionMultiplexer>().GetDatabase());

            // Add framework services.
            services.AddMvc();
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IHostingEnvironment env, ILoggerFactory loggerFactory)
        {
            loggerFactory.AddConsole(Configuration.GetSection("Logging"));
            loggerFactory.AddDebug();

            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
                app.UseBrowserLink();
            }
            else
            {
                app.UseExceptionHandler("/Home/Error");
            }

            app.UseDefaultFiles();
            app.UseStaticFiles();

            // var rewriter = new RewriteOptions();
            // rewriter.AddDomainRedirect("isitcg.herokuapp.com", "www.isitcg.com");
            // app.UseRewriter(rewriter);

            app.UseMvc(routes =>
            {
                routes.MapRoute(
                    name: "default",
                    template: "{controller=Home}/{action=Index}/{id?}");
            });
        }
    }
}