package constant

const (
	// BuildID indicates the ID for this specific service build.
	BuildID = Key("build_id")

	// CorrelationID holds the unique ID that allows to trace an execution.
	CorrelationID = Key("correlation_id")

	// InBehalfOf indicates the piece of software that launched the execution.
	// Best used along with WhosRequesting.
	InBehalfOf = Key("in_behalf_of")

	// IsDryRun can be used to determine whether a request is a dry run.
	IsDryRun = Key("is_dry_run")

	// ServiceHost indicates the host where the service is currently running.
	ServiceHost = Key("service_host")

	// ServiceName indicates the service name.
	ServiceName = Key("service_name")

	// WhosRequesting indicates the piece of software that is requesting
	// the current execution. Best used along with InBehalfOf.
	WhosRequesting = Key("whos_requesting")
)

// Key encapsulates well known data points used all over the place.
type Key string
