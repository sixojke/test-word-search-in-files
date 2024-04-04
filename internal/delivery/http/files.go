package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/sixojke/internal/domain"
)

func (h *Handler) filesSearch(c *gin.Context) {
	input := domain.FilesSearchInp{Word: c.Query("word")}

	if err := input.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})

		return
	}

	out, err := h.services.Files.FindFilesWithWord(input)
	if err != nil {
		log.Error(err)
		newErrorResponse(c, http.StatusInternalServerError, errorResponse{
			Code:    http.StatusInternalServerError,
			Message: "error.files.internalServerError",
		})

		return
	}

	c.JSON(http.StatusOK, out)
}
