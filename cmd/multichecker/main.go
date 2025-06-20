// Package multichecker объединяет множество статических анализаторов Go в единый инструмент.
//
// Включает анализаторы из нескольких источников:
//   - Стандартные анализаторы из golang.org/x/tools/go/analysis/passes
//   - Анализаторы staticcheck/honnef.co (популярные проверки качества кода)
//   - Сторонние анализаторы (ineffassign, bodyclose)
//   - Кастомные анализаторы (например, staticlint.NoOsExitAnalyzer)
//
// Особенности:
//   - Объединяет 30+ анализаторов в один исполняемый файл
//   - Поддерживает проверку всех основных аспектов кода:
//   - корректность (nilness, errorsas)
//   - производительность (loopclosure, lostcancel)
//   - безопасность (cgocall, unsafeptr)
//   - стиль кода (structtag, printf)
//   - Легко расширяется добавлением новых анализаторов
//
// Использование:
//
//	Запустите проверку на вашем коде из корня проекта:
//	  ./cmd/multichecker/multichecker ./cmd/...
package multichecker

import (
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/ramil063/gometrics/cmd/staticlint"
)

// AllAnalyzers содержит базовые анализаторы из golang.org/x/tools/go/analysis/passes.
//
// Включает следующие проверки:
//
//   - asmdecl: проверяет корректность объявлений ассемблерных файлов Go
//   - assign: обнаруживает бесполезные присваивания (x = x)
//   - atomic: проверяет корректное использование atomic-операций
//   - bools: находит подозрительные операции с булевыми значениями
//   - buildtag: проверяет корректность тегов сборки (build tags)
//   - cgocall: проверяет правила вызова C-кода через cgo
//   - composite: обнаруживает неинициализированные композитные литералы
//   - copylock: находит копирование мьютексов и других lock-объектов
//   - errorsas: проверяет корректность использования errors.As
//   - fieldalignment: предлагает оптимизацию выравнивания полей структур
//   - httpresponse: проверяет обязательную обработку HTTP-ответов
//   - ifaceassert: обнаруживает невозможные утверждения типов
//   - loopclosure: выявляет проблемные замыкания переменных в циклах
//   - lostcancel: находит утечку контекста из-за невызова cancel-функции
//   - nilfunc: обнаруживает сравнение функций с nil
//   - nilness: проверяет код на nil-паники (анализ потока данных)
//   - printf: валидирует строки форматирования в Printf-функциях
//   - shadow: находит случайное затенение переменных
//   - shift: проверяет корректность операций битового сдвига
//   - stdmethods: проверяет соответствие стандартным интерфейсам
//   - structtag: валидирует теги структур
//   - tests: обнаруживает распространенные ошибки в тестах
//   - unmarshal: проверяет обработку ошибок при анмаршалинге
//   - unreachable: находит недостижимый код
//   - unsafeptr: проверяет безопасное использование unsafe.Pointer
//   - s1000: базовые проверки простого кода
//   - st1000: проверки стиля кода
//   - qf1001: проверки с автоматическими исправлениями
//   - ineffassign: проверяет неэффективные присваивания
//   - bodyclose: проверяет не закрытые HTTP body
//   - staticlint: проверяет отсутствие вызова функции os.Exit() в main пакете main функции
//   - SA-анализаторы: продвинутые проверки корректности кода
var AllAnalyzers = []*analysis.Analyzer{
	asmdecl.Analyzer,        // asmdecl - проверяет корректность объявлений ассемблерных файлов Go
	assign.Analyzer,         // assign - обнаруживает бесполезные присваивания (x = x)
	atomic.Analyzer,         // atomic - проверяет корректное использование atomic-операций
	bools.Analyzer,          // bools - находит подозрительные операции с булевыми значениями (например, x == true вместо x)
	buildtag.Analyzer,       // buildtag - проверяет корректность тегов сборки (build tags)
	cgocall.Analyzer,        // cgocall - проверяет правила вызова C-кода через cgo
	composite.Analyzer,      // composite - обнаруживает неинициализированные композитные литералы
	copylock.Analyzer,       // copylock - находит копирование мьютексов и других lock-объектов
	errorsas.Analyzer,       // errorsas - проверяет корректность использования errors.As
	fieldalignment.Analyzer, // fieldalignment - предлагает оптимизацию выравнивания полей структур для уменьшения занимаемой памяти
	httpresponse.Analyzer,   // httpresponse - проверяет обязательную обработку HTTP-ответов
	ifaceassert.Analyzer,    // ifaceassert - обнаруживает невозможные утверждения типов (type assertions)
	loopclosure.Analyzer,    // loopclosure - выявляет проблемные замыкания переменных в циклах
	lostcancel.Analyzer,     // lostcancel - находит утечку контекста из-за невызова cancel-функции
	nilfunc.Analyzer,        // nilfunc - обнаруживает сравнение функций с nil
	nilness.Analyzer,        // nilness - проверяет код на nil-паники (анализ потока данных)
	printf.Analyzer,         // printf - валидирует строки форматирования в Printf-функциях
	shadow.Analyzer,         // shadow - находит случайное затенение переменных
	shift.Analyzer,          // shift - проверяет корректность операций битового сдвига
	stdmethods.Analyzer,     // stdmethods - проверяет соответствие стандартным интерфейсам (например, String() string)
	structtag.Analyzer,      // structtag - валидирует теги структур (например, `json:"name"`)
	tests.Analyzer,          // tests - обнаруживает распространенные ошибки в тестах
	unmarshal.Analyzer,      // unmarshal - проверяет обработку ошибок при анмаршалинге
	unreachable.Analyzer,    // unreachable - находит недостижимый код
	unsafeptr.Analyzer,      // unsafeptr - проверяет безопасное использование unsafe.Pointer
}

// Main является точкой входа для мультианализатора.
//
// Собирает все анализаторы в один набор и запускает их через multichecker.Main.
// Анализаторы разделены на несколько групп:
//  1. Базовые анализаторы из x/tools/go/analysis/passes
//  2. Анализаторы staticcheck (все доступные проверки)
//  3. Дополнительные сторонние анализаторы:
//     - ineffassign - обнаружение неиспользуемых присваиваний
//     - bodyclose - проверка закрытия HTTP response bodies
//  4. Кастомные анализаторы (например, запрет os.Exit в main)
func Main() {

	// Добавляем дополнительные анализаторы
	AllAnalyzers = append(AllAnalyzers, getStaticChecks()...)
	AllAnalyzers = append(
		AllAnalyzers,
		simple.Analyzers[0].Analyzer,     // Базовые проверки простого кода
		stylecheck.Analyzers[0].Analyzer, // Проверки стиля кода
		quickfix.Analyzers[0].Analyzer,   // Проверки с автоматическими исправлениями
		ineffassign.Analyzer,             // Неэффективные присваивания
		bodyclose.Analyzer,               // Не закрытые HTTP body
		staticlint.NoOsExitAnalyzer,      // Вызов функции os.Exit() в main пакете main функции
	)

	// Передаем как variadic аргумент
	multichecker.Main(AllAnalyzers...)
}

// getStaticChecks возвращает все анализаторы из staticcheck.
//
// Staticcheck предоставляет сотни специализированных проверок,
// сгруппированных по категориям. Эта функция включает их все.
// Возвращаемый слайс содержит анализаторы для:
//   - Проверки корректности (SA*)
//   - Проверки производительности (perf*)
//   - Проверки стиля (ST*)
//   - Других категорий
func getStaticChecks() []*analysis.Analyzer {
	var myChecks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}
	return myChecks
}
