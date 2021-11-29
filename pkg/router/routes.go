package router

import (
	"MdShorts/pkg/account_service"
	"MdShorts/pkg/bookmark_service"
	"MdShorts/pkg/category_service"
	"MdShorts/pkg/news_service"
	"MdShorts/pkg/profile_service"
	"MdShorts/pkg/share_service"
	"MdShorts/pkg/unregistered_user_service"
	user_news_check_service "MdShorts/pkg/userNewsCheck_service"
	"MdShorts/pkg/util_service"
	"net/http"
)

// Route type description
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes contains all routes
type Routes []Route

var routes = Routes{
	Route{
		"signup",
		"POST",
		"/signup",
		account_service.SignUp,
	},
	Route{
		"login",
		"POST",
		"/login",
		account_service.Login,
	},
	Route{
		"verifyemail",
		"GET",
		"/verifyemail/{code}",
		account_service.VerifyEmail,
	},
	Route{
		"verifymnumb",
		"POST",
		"/verifymobilenumber",
		account_service.VerifyOTP,
	},
	Route{
		"resendOtp",
		"POST",
		"/resendotp",
		account_service.ResendOTP,
	},
	Route{
		"sendemail",
		"POST",
		"/sendemail",
		account_service.SendVerificationEmail,
	},
	Route{
		"sendemail",
		"GET",
		"/resetpassword",
		account_service.SendPasswordResetLink,
	},
	Route{
		"update profile",
		"PUT",
		"/profile",
		profile_service.UpdateProfile,
	},
	Route{
		"update profile",
		"PUT",
		"/password/profile",
		profile_service.ChangePassword,
	},
	Route{
		"set profile",
		"POST",
		"/profile",
		profile_service.SetProfile,
	},
	Route{
		"get profile",
		"GET",
		"/profile",
		profile_service.Profile,
	},
	Route{
		"getCategory",
		"GET",
		"/category",
		category_service.GetAllCategory,
	},
	Route{
		"updateCategory",
		"PUT",
		"/category/{categoryId}",
		category_service.UpdateCategory,
	},
	Route{
		"getCategoryByIds",
		"GET",
		"/category/{categoryIds}",
		category_service.GetCategoriesWithIDs,
	},
	Route{
		"category",
		"POST",
		"/category",
		category_service.AddCategory,
	},
	Route{
		"utilPreSigned",
		"POST",
		"/util/presignedurl",
		util_service.GetPreSignedURL,
	},
	Route{
		"getnews",
		"GET",
		"/news/{userId}",
		news_service.GetNews,
	},
	Route{
		"getGlobalnews",
		"GET",
		"/gnews",
		news_service.GetGlobalNews,
	},
	Route{
		"add news status for user",
		"POST",
		"/addnews",
		user_news_check_service.AddUserNewsCheck,
	},
	Route{
		"update news status for user",
		"PUT",
		"/updatenews",
		user_news_check_service.UpdateUserNewsCheck,
	},
	Route{
		"getnews",
		"GET",
		"/news/",
		news_service.GetGlobalNews,
	},
	Route{
		"getnews",
		"GET",
		"/search/news",
		news_service.GetSearchNews,
	},
	Route{
		"share",
		"POST",
		"/share",
		share_service.AddShare,
	},
	Route{
		"share",
		"GET",
		"/share",
		share_service.GetShares,
	},
	Route{
		"bookmark",
		"POST",
		"/bookmark",
		bookmark_service.AddBookmark,
	},
	Route{
		"bookmark",
		"GET",
		"/bookmark",
		bookmark_service.GetBookmarks,
	},
	Route{
		"bookmark",
		"PUT",
		"/bookmark/{bookmarkId}",
		bookmark_service.UpdateBookmark,
	},
	Route{
		"unregisteruser",
		"POST",
		"/add/unregisteruser",
		unregistered_user_service.AddUnregisteredUserService,
	},
	Route{
		"unregisteruser",
		"GET",
		"/unregisteruser",
		unregistered_user_service.AddUnregisteredUserService,
	},
}
