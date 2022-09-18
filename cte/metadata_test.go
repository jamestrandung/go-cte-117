package cte

import (
    "reflect"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestExtractMetadata(t *testing.T) {
    defer func(original func(t reflect.Type) reflect.Type) {
        extractNonPointerType = original
    }(extractNonPointerType)

    scenarios := []struct {
        desc string
        test func(t *testing.T)
    }{
        {
            desc: "nil metadata",
            test: func(t *testing.T) {
                mpMock := &MockMetadataProvider{}

                mpMock.On("CTEMetadata").
                    Return(nil).
                    Twice()

                assert.PanicsWithError(
                    t, ErrNilMetadata.Err(reflect.TypeOf(mpMock)).Error(), func() {
                        extractMetadata(mpMock, true)
                    },
                )

                assert.PanicsWithError(
                    t, ErrNilMetadata.Err(reflect.TypeOf(mpMock)).Error(), func() {
                        extractMetadata(mpMock, false)
                    },
                )

                mock.AssertExpectationsForObjects(t, mpMock)
            },
        },
        {
            desc: "isComputerKey == false",
            test: func(test *testing.T) {
                mpMock := &MockMetadataProvider{}

                dummyMetadata := "something"

                mpMock.On("CTEMetadata").
                    Return(dummyMetadata).
                    Once()

                extractNonPointerType = func(t reflect.Type) reflect.Type {
                    assert.Equal(test, reflect.TypeOf(dummyMetadata), t)

                    return reflect.TypeOf(
                        struct {
                            field1 string
                            field2 int
                        }{},
                    )
                }

                metadata := extractMetadata(mpMock, false)
                assert.Equal(test, 2, len(metadata))
                assert.Equal(test, reflect.TypeOf("string"), metadata["field1"])
                assert.Equal(test, reflect.TypeOf(1), metadata["field2"])

                mock.AssertExpectationsForObjects(t, mpMock)
            },
        },
        {
            desc: "isComputerKey == true",
            test: func(test *testing.T) {
                mpMock := &MockMetadataProvider{}

                dummyMetadata := "something"

                mpMock.On("CTEMetadata").
                    Return(dummyMetadata).
                    Once()

                extractNonPointerType = func(t reflect.Type) reflect.Type {
                    assert.Equal(test, reflect.TypeOf(dummyMetadata), t)

                    return reflect.TypeOf(
                        struct {
                            field1 string
                            field2 int
                        }{},
                    )
                }

                metadata := extractMetadata(mpMock, true)
                assert.Equal(test, 3, len(metadata))
                assert.Equal(test, reflect.TypeOf("string"), metadata["field1"])
                assert.Equal(test, reflect.TypeOf(1), metadata["field2"])
                assert.Equal(test, reflect.TypeOf(mpMock), metadata[metaTypeComputerKey])

                mock.AssertExpectationsForObjects(t, mpMock)
            },
        },
    }

    for _, scenario := range scenarios {
        s := scenario

        t.Run(s.desc, s.test)
    }
}

func TestParsedMetadata_GetComputerKeyType(t *testing.T) {
    var pm parsedMetadata = make(map[metaType]reflect.Type)

    result, ok := pm.getComputerKeyType()
    assert.Equal(t, reflect.Type(nil), result)
    assert.False(t, ok)

    pm[metaTypeComputerKey] = reflect.TypeOf("dummy")

    result, ok = pm.getComputerKeyType()
    assert.Equal(t, reflect.TypeOf("dummy"), result)
    assert.True(t, ok)
}

func TestParsedMetadata_GetComputerType(t *testing.T) {
    var pm parsedMetadata = make(map[metaType]reflect.Type)

    result, ok := pm.getComputerType()
    assert.Equal(t, reflect.Type(nil), result)
    assert.False(t, ok)

    pm[metaTypeComputer] = reflect.TypeOf("dummy")

    result, ok = pm.getComputerType()
    assert.Equal(t, reflect.TypeOf("dummy"), result)
    assert.True(t, ok)
}

func TestParsedMetadata_GetInoutInterface(t *testing.T) {
    var pm parsedMetadata = make(map[metaType]reflect.Type)

    result, ok := pm.getInoutInterface()
    assert.Equal(t, reflect.Type(nil), result)
    assert.False(t, ok)

    pm[metaTypeInout] = reflect.TypeOf("dummy")

    result, ok = pm.getInoutInterface()
    assert.Equal(t, reflect.TypeOf("dummy"), result)
    assert.True(t, ok)
}
