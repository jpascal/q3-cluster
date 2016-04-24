package uuid_test

import (
	"testing"
	"uuid"
	"github.com/stretchr/testify/assert"
	"regexp"
	"fmt"
)

func TestNewUUID(t *testing.T) {
	id := uuid.NewUUID()
	n, _ := regexp.Match("[a-f0-9]{8}(?:-[a-f0-9]{4}){4}-[a-f0-9]{8}", []byte(id))
	assert.Equal(t, n, true, fmt.Sprintf("'%v' is not UUID",id))
}
