package interfaces

type ICommand interface {
	Run(flags IFlag)
}
