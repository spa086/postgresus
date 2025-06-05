package notifiers

type NotificationSender interface {
	Send(heading string, message string) error

	Validate() error
}
