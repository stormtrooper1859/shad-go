# distbuild

В этом задании вам нужно будет реализовать систему распределённой сборки.

Система сборки получает на вход граф сборки и файлы с исходным кодом. Результатом сборки
являются исполняемые файлы и stderr/stdout запущенных процессов.

## Граф сборки

Граф сборки состоит из джобов. Каждый джоб описывает команды, которые нужно запустить на одной машине,
вместе со всеми входными файлами, которые нужны этим командам для работы.

Джобы в графе сборки запускают произвольные команды. Например, вызывать компилятор, линкер или 
запускать тесты.

Команды внутри джоба могут читать файлы с файловой системы. Мы будем различать два вида файлов:
 - Файлы с исходным кодом с машины пользователя.
 - Файлы, которые породили другие джобы.

Команды внутри джоба могут писать результаты своей работы в файлы на диске. Выходные файлы
обязаны находиться внутри его выходной директории. Директория с результатом работы джоба называется
артефактом.

```go
package build

import "crypto/sha1"

// ID задаёт уникальный идентификатор джоба.
//
// Мы будем использовать sha1 хеш, поэтому ID будет занимать 20 байт.
type ID [sha1.Size]byte

// Job описывает одну вершину графа сборки.
type Job struct {
	// ID задаёт уникальный идентификатор джоба.
	//
	// ID вычисляется как хеш от всех входных файлов, команд запуска и хешей зависимых джобов.
	//
	// Выход джоба целиком определяется его ID. Это важное свойство позволяет кешировать
	// результаты сборки.
	ID ID

	// Name задаёт человекочитаемое имя джоба.
	//
	// Например:
	//   build gitlab.com/slon/disbuild/pkg/b
	//   vet gitlab.com/slon/disbuild/pkg/a
	//   test gitlab.com/slon/disbuild/pkg/test
	Name string

	// Inputs задаёт список файлов из директории с исходным кодом,
	// которые нужны для работы этого джоба.
	//
	// В типичном случае, тут будут перечислены все .go файлы одного пакета.
	Inputs []string

	// Deps задаёт список джобов, выходы которых нужны для работы этого джоба.
	Deps []ID

	// Cmds описывает список команд, которые нужно выполнить в рамках этого джоба.
	Cmds []Cmd
}
```

## Архитектура системы

Наша система будет состоять из трех компонент.
 * Клиент - процесс запускающий сборку.
 * Воркер - процесс запускающий команды компиляции и тестирования.
 * Координатор - центральный процесс в системе, общается с клиентами и воркерами. Раздаёт задачи
   воркерам.

Типичная сборка выглядит так:
1. Клиент подключается к координатору, посылает ему граф сборки и входные файлы для графа сборки.
2. Координатор сохраняет граф сборки в памяти и начинает его исполнение.
3. Воркеры начинают выполнять вершины графа, пересылая друг другу выходные директории джобов.
4. Результаты работы джобов скачиваются на клиента.

# Как решать эту задачу

Задача разбита на шаги. В начале, вам нужно будет реализовать небольшой набор независимых пакетов,
которые реализует нужные примитивы. Код в этих пакетах покрыт юниттестами. В каждом пакете находится
файл README.md, объясняющий подзадачу.

Рекомендуемый порядок выполнения:

- [`distbuild/pkg/build`](./pkg/build) - определение графа сборки. В этом пакете ничего писать не нужно,
  нужно ознакомиться с существующим кодом.
- [`distbuild/pkg/tarstream`](./pkg/tarstream) - передача директории через сокет.
- [`distbuild/pkg/api`](./pkg/api) - протокол общения между компонентами.
- [`distbuild/pkg/artifact`](./pkg/artifact) - кеш артефактов и протокол передачи артефактов между воркерами.
- [`distbuild/pkg/filecache`](./pkg/filecache) - кеш файлов и протокол передачи файлов между компонентами.
- [`distbuild/pkg/scheduler`](./pkg/scheduler) - планировщик с эвристикой локальности.

После того, как все кубики будут готовы, нужно будет соединить их вместе, реализовав `distbuild/pkg/worker`,
`distbuild/pkg/client` и `distbuild/pkg/dist`. Код в этих пакетах нужно отлаживать на
интеграционных тестах в [`distbuild/disttest`](../disttest).

Код тестов в этом задании менять нельзя. Это значит, что вы не можете менять интерфейсы в тех местах, где
код покрыт тестами.

<details>
  <summary markdown="span">Сколько кода нужно написать?</summary>
  
  ```
prime@bee ~/C/shad-go> find distbuild -iname '*.go' | grep -v test | grep -v mock | grep -v pkg/build | xargs wc -l
   23 distbuild/pkg/worker/state.go
  111 distbuild/pkg/worker/worker.go
   45 distbuild/pkg/worker/download.go
  281 distbuild/pkg/worker/job.go
   69 distbuild/pkg/api/heartbeat.go
  121 distbuild/pkg/api/build_client.go
   53 distbuild/pkg/api/build.go
   60 distbuild/pkg/api/heartbeat_handler.go
  142 distbuild/pkg/api/build_handler.go
   56 distbuild/pkg/api/heartbeat_client.go
  288 distbuild/pkg/scheduler/scheduler.go
  119 distbuild/pkg/dist/build.go
  120 distbuild/pkg/dist/coordinator.go
   98 distbuild/pkg/tarstream/stream.go
   42 distbuild/pkg/artifact/client.go
  191 distbuild/pkg/artifact/cache.go
   54 distbuild/pkg/artifact/handler.go
  124 distbuild/pkg/client/build.go
   83 distbuild/pkg/filecache/client.go
   99 distbuild/pkg/filecache/handler.go
  111 distbuild/pkg/filecache/filecache.go
 2290 total
  ```
</details>
