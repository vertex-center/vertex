package app

type Meta struct {
	// ID is the unique identifier of the app.
	ID string `json:"id"`

	// Name is the name of the app visible to the user.
	Name string `json:"name"`

	// Description is a brief description of the app.
	Description string `json:"description"`

	// Icon is the material symbol name for the app.
	Icon string `json:"icon"`

	// Category is the category of the app.
	Category string `json:"category"`

	// DefaultPort is the default port of the app.
	DefaultPort string `json:"port"`

	// DefaultKernelPort is the default port of the app in kernel mode.
	DefaultKernelPort string `json:"kernel_port"`

	// Hidden is a flag that indicates if the app does only backend work and should be hidden from the frontend.
	Hidden bool `json:"hidden"`

	// Dependencies is a list of app IDs that this app depends on.
	Dependencies []*Meta `json:"-"`
}
