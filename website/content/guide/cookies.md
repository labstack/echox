+++
title = "Cookies"
description = "Handling cookie in Echo"
[menu.main]
  name = "Cookies"
  parent = "guide"
+++

Cookie is a small piece of data sent from a website server and stored in the user's web
browser while browsing. Every time the user loads the website, the browser
sends the cookies back to the server to notify the server of user's latest activity.
Cookies were designed to be a reliable mechanism for websites to remember stateful
information (e.g. items added to the shopping cart in an online store) or to
record the user's browsing activity (such as clicking particular buttons, logging
in, or user previously visited pages of the website). Cookies can also store form content a user has previously entered, such as username, gender, age, address, etc.

## Cookie Attributes

Attribute | Optional
:--- | :---
`Name` | No
`Value` | No
`Path` | Yes
`Domain` | Yes
`Expires` | Yes
`Secure` | Yes
`HttpOnly` | Yes

Echo uses go standard `http.Cookie` object to add/retrieve cookies from the context received in the handler function.

## Create a Cookie

```go
func writeCookie(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = "jon"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
	return c.String(http.StatusOK, "write a cookie")
}
```

- Cookie is created using `new(http.Cookie)`.
- Attributes for the cookie are set assigning to the `http.Cookie` instance public attributes.  
- Finally `c.SetCookie(cookie)` adds a `Set-Cookie` header in HTTP response.

## Read a Cookie

```go
func readCookie(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)
	return c.String(http.StatusOK, "read a cookie")
}
```

- Cookie is read by name using `c.Cookie("username")` from the HTTP request.
- Cookie attributes are accessed using `Getter` function.

## Read all the Cookies

```go
func readAllCookies(c echo.Context) error {
	for _, cookie := range c.Cookies() {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
	}
	return c.String(http.StatusOK, "read all the cookies")
}
```
