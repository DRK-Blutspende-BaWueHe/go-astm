package e2e

import (
	"fmt"
	"github.com/DRK-Blutspende-BaWueHe/go-astm/astm1384"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadfileBeCom52(t *testing.T) {
	fileData, err := ioutil.ReadFile("../protocoltest/becom/5.2/bloodtype.astm")
	if err != nil {
		fmt.Println("Failed : ", err)
		return
	}

	message, err := astm1384.Unmarshal(fileData,
		astm1384.EncodingWindows1252, astm1384.TimezoneEuropeBerlin, astm1384.LIS2A2)
	if err != nil {
		fmt.Println("Error in unmarshaling the message ", err)
		return
	}

	// check header
	assert.Equal(t, "IH v5.2", message.Header.SenderStreetAddress)
	assert.Equal(t, "Bio-Rad", message.Header.SenderNameOrID)

	// Message was recorded in Germany (UTC+1) 19:42
	locale, err := time.LoadLocation("Europe/Berlin")
	assert.Nil(t, err)
	localtime := message.Header.DateAndTime.In(locale)
	assert.Equal(t, "20220315194227", localtime.Format("20060102150405"))

	assert.Equal(t, 1 /*Patient*/, len(message.Records))
	assert.Equal(t, "1010868845", message.Records[0].Patient.LabAssignedPatientID)
	assert.Equal(t, "Testus", message.Records[0].Patient.LastName)
	assert.Equal(t, "Test", message.Records[0].Patient.FirstName)
	assert.Equal(t, "19400607", message.Records[0].Patient.DOB.Format("20060102"))
	assert.Equal(t, "M", message.Records[0].Patient.Gender)

	// Check for the results
	assert.Equal(t, 1, message.Records[0].OrdersAndResults[0].Order.SequenceNumber)
	assert.Equal(t, "1122206642", message.Records[0].OrdersAndResults[0].Order.SpecimenID[0])
	//assert.Equal(t, "1122206642", message.Records[0].Orders[0].Order.InstrumentSpecimenID[0])
	//assert.Equal(t, "1122206642", message.Records[0].Orders[0].Order.InstrumentSpecimenID[4])
	// this is an empty field with a date (time.IsZero)
	assert.Equal(t, "", message.Records[0].OrdersAndResults[0].Order.CollectionID)
	// assert.Equal(t, 7 /*Fields to be scanned*/, len(message.Records[0].Orders[0].Order.InstrumentSpecimenID))
	assert.Equal(t, "MO10", message.Records[0].OrdersAndResults[0].Order.UniversalTestID_ManufacturerCode)
	assert.Equal(t, "28343" /*Lot#*/, message.Records[0].OrdersAndResults[0].Order.UniversalTestID_Custom2)
	assert.Equal(t, "R" /*Routine*/, message.Records[0].OrdersAndResults[0].Order.Priority)
	assert.Equal(t, 7 /*Fields to be scanned*/, len(message.Records[0].OrdersAndResults[0].Order.UniversalTestID))

	assert.Equal(t, "20220311093217" /*UTC Time (this is -1 to what is observed in file)*/, message.Records[0].OrdersAndResults[0].Order.RequestedOrderDateTime.Format("20060102150405"))
	assert.Equal(t, "20220311093217", message.Records[0].OrdersAndResults[0].Order.SpecimenCollectionDateTime.Format("20060102150405"))
	assert.Equal(t, "11", message.Records[0].OrdersAndResults[0].Order.UserField1)
	assert.Equal(t, "20220311104103", message.Records[0].OrdersAndResults[0].Order.DateTimeResultsReported.Format("20060102150405"))
	assert.Equal(t, "P" /*Preliminary*/, message.Records[0].OrdersAndResults[0].Order.ReportType)

	/* fmt.Printf("Messageheader: %+v\n", message.Header)
	   for _, record := range message.Records {
		fmt.Printf("Patient : %s, %s\n", record.Patient.Name[0], record.Patient.Name[1])
		for _, order := range record.Orders {
			fmt.Printf("  Order: %+v\n", order.Order)
			for _, result := range order.Results {
				fmt.Printf("   Result: %+v\n", result.Result)
			}
		}
	}*/
}