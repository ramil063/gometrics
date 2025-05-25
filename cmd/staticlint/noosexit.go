// Package staticlint содержит набор статических анализаторов для проверки кода.
//
// Анализаторы включают:
//   - Проверку запрещённых вызовов os.Exit в main-пакете
//   - [другие анализаторы, если есть]
//
// Особенности:
//   - Интеграция с analysis package (golang.org/x/tools/go/analysis)
//   - Поддержка мультичекеров (multichecker)
//   - Фокусировка на соглашениях и best practices
package staticlint

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NoOsExitAnalyzer проверяет отсутствие вызовов os.Exit() в main-пакете.
//
// Такой вызов может помешать выполнению отложенных функций (defer)
// и корректному завершению приложения. Анализатор сообщает о каждом
// обнаруженном вызове os.Exit() в main-пакете как об ошибке.
var NoOsExitAnalyzer = &analysis.Analyzer{
	Name:     "noosexit",
	Doc:      "forbid direct os.Exit calls in main package",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// run реализует основную логику анализатора NoOsExitAnalyzer.
//
// Функция:
//  1. Проверяет, что анализ выполняется для main-пакета(игнорирует файлы кеша Go)
//  2. Ищет все вызовы функций в коде
//  3. Фильтрует вызовы os.Exit()
//  4. Сообщает о найденных случаях через pass.Reportf()
//
// Возвращает (nil, nil) если анализ выполнен без ошибок,
// даже если найдены нарушения (объявляются через Reportf).
func run(pass *analysis.Pass) (interface{}, error) {
	// Проверяем, что мы в пакете main
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	// Игнорируем файлы в кеше Go
	for _, f := range pass.Files {
		pos := pass.Fset.Position(f.Pos())
		if strings.Contains(pos.Filename, "/.cache/go-build/") ||
			strings.Contains(pos.Filename, "/go/pkg/mod/") {
			return nil, nil
		}
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Фильтр для функций main
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)

		// Проверяем, что это функция main
		if fn.Name.Name != "main" {
			return
		}

		// Ищем вызовы os.Exit
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if x, ok := sel.X.(*ast.Ident); ok {
				if x.Name == "os" && sel.Sel.Name == "Exit" {
					pass.Reportf(call.Pos(),
						"direct os.Exit call in main function is forbidden")
				}
			}
			return true
		})
	})

	return nil, nil
}
