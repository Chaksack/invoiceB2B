package middleware

import (
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories" // To fetch staff role
	"invoiceB2B/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AdminMiddleware struct {
	staffRepo repositories.StaffRepository
}

func NewAdminMiddleware(staffRepo repositories.StaffRepository) *AdminMiddleware {
	return &AdminMiddleware{staffRepo: staffRepo}
}

// AdminRequired checks if the authenticated user is an admin.
// It assumes that the general AuthMiddleware.Protected() has already run
// and validated the JWT, placing user claims in c.Locals("user").
// It also assumes that admin users are stored in the 'staff' table.
func (am *AdminMiddleware) AdminRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals("user").(*jwt.Token)
		if !ok {
			// This should not happen if AuthMiddleware.Protected() ran correctly
			return utils.HandleError(c, fiber.StatusUnauthorized, "Authentication token not found in context.", nil)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
		}

		// Check if the token subject indicates it's for a staff member or if there's a specific role claim.
		// For this example, we assume the user_id in the JWT for an admin corresponds to an ID in the 'staff' table.
		// A more robust way is to have a 'role' claim in the JWT.

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return utils.HandleError(c, fiber.StatusForbidden, "User ID not found in token claims.", nil)
		}

		staffID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusForbidden, "Invalid user ID format in token for admin check.", err)
		}

		// Fetch staff member by ID to check their role
		staff, err := am.staffRepo.FindByID(c.Context(), uint(staffID))
		if err != nil || staff == nil {
			// If not found in staff table, or DB error, then not an authorized admin
			return utils.HandleError(c, fiber.StatusForbidden, "Access denied. Admin privileges required.", nil)
		}

		// Check the role (e.g., staff.Role == "admin")
		// For simplicity, any entry in the staff table is considered having some level of admin access.
		// Refine this with specific role checks (e.g., models.RoleAdmin)
		if staff.Role != models.RoleAdmin && staff.Role != models.RoleFinanceManager && staff.Role != models.RoleKYCReviewer { // Example roles
			// This check depends on how you define admin roles.
			// If any staff is admin, this check might be simpler.
			// If specific roles are needed for specific admin routes, this middleware might need to be more granular
			// or you'd have multiple admin middlewares.
			// For now, let's assume certain roles are "admin-level".
			// A common approach is to have a general "admin" role.
			// If the JWT itself contained a "role": "admin" claim, this DB lookup could be skipped.
			// For now, we check if the staff member has a role that grants admin access.
			// This is a placeholder for more specific role checking.
			// For this example, we'll just check if the role is "admin".
			if staff.Role != models.RoleAdmin {
				return utils.HandleError(c, fiber.StatusForbidden, "Access denied. Insufficient admin privileges.", nil)
			}
		}

		// Store staff details in context if needed by admin handlers
		c.Locals("staff_id", staff.ID)
		c.Locals("staff_role", staff.Role)

		return c.Next()
	}
}
