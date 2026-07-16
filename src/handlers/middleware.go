package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// SecurityHeaders adds security headers to all responses.
func SecurityHeaders(c fiber.Ctx) error {
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Frame-Options", "DENY")
	c.Set("Content-Security-Policy",
		"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline'; "+
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://cdnjs.cloudflare.com; "+
			"font-src 'self' https://fonts.gstatic.com https://cdnjs.cloudflare.com; "+
			"img-src 'self' data: https://copilot.microsoft.com https://images.unsplash.com https://thumbs.dreamstime.com; "+
			"connect-src 'self';")
	c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
	return c.Next()
}

// AuthRequired ensures a valid session exists.
func AuthRequired(c fiber.Ctx) error {
	sess := session.FromContext(c)
	userID, ok := sess.Get("user_id").(int)
	if !ok || userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "No autenticado"})
	}
	c.Locals("user_id", userID)
	c.Locals("role", sess.Get("role"))
	return c.Next()
}

// AdminRequired ensures the user has admin role.
func AdminRequired(c fiber.Ctx) error {
	sess := session.FromContext(c)
	role, ok := sess.Get("role").(string)
	if !ok || role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acceso denegado"})
	}
	c.Locals("role", "admin")
	return c.Next()
}
