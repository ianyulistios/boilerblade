package cli

import "fmt"

func ShowHelp() {
	fmt.Println("Boilerblade CLI - Go Boilerplate Generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  boilerblade <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  new <project-name>     Create a new Boilerblade project")
	fmt.Println("  make <resource>        Generate code (model, repository, usecase, handler, dto, consumer, migration, all)")
	fmt.Println("  version                Show version information")
	fmt.Println("  help                    Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  boilerblade new my-api")
	fmt.Println("  boilerblade make model Product")
	fmt.Println("  boilerblade make all Product -fields=\"Name:string:required,Price:float64:required\"")
	fmt.Println("  boilerblade make consumer -name=OrderEvents -title=\"Order Events\"")
	fmt.Println("  boilerblade make migration -name=add_orders_table")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/ianyulistio/boilerblade")
}
