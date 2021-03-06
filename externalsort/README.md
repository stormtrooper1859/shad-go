## externalsort

В этой задаче нужно написать однопроходную внешнюю сортировку слиянием.
Моделируется ситуация, в которой данные расположены на внешних устройствах и суммарно не вмещаются в оперативную память,
но каждый кусочек по отдельности вмещается.

Задача разбита на 3 составные части.

#### Reader & writer

Реализовать интерфейсы для построчного чтения/записи строк:
```
type LineReader interface {
	ReadLine() (string, error)
}

type LineWriter interface {
	Write(l string) error
}
```
и два конструктора:
```
func NewReader(r io.Reader) LineReader
func NewWriter(w io.Writer) LineWriter
```

`NewLineReader` оборачивает переданный `io.Reader` в `LineReader`.

Вызов `ReadLine` должен читать одну строку.
Строка имеет произвольную длину.
Конец строки определяется переводом строки ('\n').
Непустая последовательность символов после последнего перевода строки также считается строкой.

`ReadLine` должен возвращать `io.EOF` при достижении конца файла.

#### Merge

Функция слияния произвольного количества отсортированных групп строк:
```
func Merge(w LineWriter, readers ...LineReader) error
```

`Merge` по необходимости читает из reader'ов и пишет во writer.

#### Sort

```
Sort(w io.Writer, in ...string) error
```

Функция принимает на вход произвольное количество файлов, каждый из которых помещается в оперативную память,
а также writer для записи результата.

Результаты сортировки отдельных файлов можно записывать поверх входных данных.

### Ссылки

* container/heap: https://golang.org/pkg/container/heap/
