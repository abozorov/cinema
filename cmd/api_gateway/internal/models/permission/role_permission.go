package permission

var RolePermission = map[string]map[string]struct{}{
	Admin: {
		MovieCreate: {},
		MovieUpdate: {},
		MovieDelete: {},
		MovieView:   {},
		MovieList:   {},

		BookingList:   {},
		BookingView:   {},
		BookingCancel: {},

		UserView: {},
		UserList: {},
	},

	User: {
		MovieList: {},
		MovieView: {},

		BookingCreate: {},
		BookingViewMe: {},
		BookingCancel: {},

		UserViewMe: {},
		UserUpdate: {},
	},
}
