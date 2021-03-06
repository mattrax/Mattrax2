package enrollpolicy

import (
	"log"
	"net/http"
	"strconv"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/generic"
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
	"github.com/pkg/errors"
)

func Handler(server *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode request from client
		var cmd Request
		if err := xml.NewDecoder(r.Body).Decode(&cmd); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest) // TODO: Correct Spec way of doing this??
			return
		}

		// Verify request
		if err := cmd.Verify(server.Config, server.UserService); err != nil {
			log.Println(errors.Wrap(err, "invalid MdePoliciesRequest:"))
			w.WriteHeader(http.StatusBadRequest) // TODO: Correct Spec way of doing this??
			return
		}

		// TODO: Use the Request Body to work out if changes have happened for PolciiesNotChanged. What is requestFilter???
		// TODO: This response is hardcoded. Automatically generate from CertService + Settings

		res := ResponseEnvelope{
			NamespaceS: "http://www.w3.org/2003/05/soap-envelope",
			NamespaceA: "http://www.w3.org/2005/08/addressing",
			NamespaceU: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd",
			HeaderAction: soap.MustUnderstand{
				MustUnderstand: "1",
				Value:          "http://schemas.microsoft.com/windows/pki/2009/01/enrollmentpolicy/IPolicy/GetPoliciesResponse",
			},
			// HeaderActivityID: GenerateActivityID(),
			HeaderRelatesTo: cmd.Header.MessageID,
			Body: ResponseBody{
				NamespaceXSI: "http://www.w3.org/2001/XMLSchema-instance",
				NamespaceXSD: "http://www.w3.org/2001/XMLSchema",
				PoliciesResponse: Response{
					PolicyID:           generic.GenerateID(), // TODO: Does this have to stay constant????
					PolicyFriendlyName: "Mattrax Identity",   // TODO: Does it show
					NextUpdateHours:    12,                   // TODO: After 12 hours does it request this same endpoint like Apple?????????
					PoliciesNotChanged: false,                // TODO: Track this. False means policies have changed since last updateHour
					Policies: []MdePolicy{
						MdePolicy{
							OIDReference: 0,
							CAs: MdeCACollection{
								Nil: true,
							},
							Attributes: MdeAttributes{
								CommonName:   "Mattrax Identity2", // TODO: Does it show
								PolicySchema: 3,
								CertificateValidity: MdeAttributesCertificateValidity{
									// TODO: What is good for these values. Also how does renewal work??
									ValidityPeriodSeconds: 1209600,
									RenewalPeriodSeconds:  172800,
								},
								EnrollmentPermission: MdeEnrollmentPermission{
									Enroll:     true,  // TODO: Try false as rejection
									AutoEnroll: false, // TODO: See what is changes
								},
								PrivateKeyAttributes: MdePrivateKeyAttributes{
									MinimalKeyLength: 2048, // TODO: Get from CertService
									KeySpec: MdeKeySpec{
										Nil: true,
									},
									KeyUsageProperty: MdeKeyUsageProperty{
										Nil: true,
									},
									Permissions: MdePermissions{
										Nil: true,
									},
									AlgorithmOIDReference: MdeAlgorithmOIDReference{
										Nil: true,
									},
									CryptoProviders: MdeCryptoProviders{
										Nil: true,
									},
								},
								Revision: MdeRevision{
									MajorRevision: 101, // TODO: Change and see what happens. Version control inside Mattrax???
									MinorRevision: 0,   // TODO: Change and see what happens.
								},
								SupersededPolicies: MdeSupersededPolicies{
									Nil: true,
								},
								PrivateKeyFlags: MdePrivateKeyFlags{
									Nil: true,
								},
								SubjectNameFlags: MdeSubjectNameFlags{
									Nil: true,
								},
								EnrollmentFlags: MdeEnrollmentFlags{
									Nil: true,
								},
								GeneralFlags: MdeGeneralFlags{
									Nil: true,
								},
								HashAlgorithmOIDReference: 0, // TODO: What dis do?
								RARequirements: MdeRARequirements{
									Nil: true,
								},
								KeyArchivalAttributes: MdeKeyArchivalAttributes{
									Nil: true,
								},
								Extensions: MdeExtensions{
									Nil: true,
								},
							},
						},
					},
				},
			},
		}

		// Marshal and send the response to client
		if response, err := xml.Marshal(res); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			/* TEMP */
			// response = []byte(`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing"><s:Body><s:Fault><s:Code><s:Value>s:Sender</s:Value><s:Subcode><s:Value xmlns:a="http://schemas.microsoft.com/net/2005/12/windowscommunicationfoundation/dispatcher">a:InternalServiceFault</s:Value></s:Subcode></s:Code><s:Reason><s:Text xml:lang="en-US">InvalidMessage</s:Text></s:Reason><s:Detail><DeviceEnrollmentServiceError xmlns="http://schemas.microsoft.com/windows/pki/2009/01/enrollment"><ErrorType>EnrollmentServer</ErrorType><Message>EnrollmentInternalServiceError</Message><TraceId>efec7f6f-7d54-4430-b56c-aeb41685d601</TraceId></DeviceEnrollmentServiceError></s:Detail></s:Fault></s:Body></s:Envelope>`)
			// w.WriteHeader(http.StatusInternalServerError)
			// 			response = []byte(`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing">
			// 	<s:Header>
			// 		<a:Action s:mustunderstand="1">http://schemas.microsoft.com/windows/pki/2009/01/enrollmentpolicy/IPolicy/GetPoliciesResponse</a:Action>
			// 		<ActivityID xmlns="http://schemas.microsoft.com/2004/09/servicemodel/diagnostics">` + generic.GenerateID() + `</ActivityID>
			// 		<a:RelatesTo>` + cmd.Header.Action + `</a:RelatesTo>
			// 	</s:Header>
			// 	<s:Body>
			// 		<s:Fault>
			// 			<s:Code>
			// 				<s:Value>s:receiver</s:value>
			// 				<s:Subcode>
			// 					<s:Value>s:Authorization</s:Value>
			// 				</s:Subcode>
			// 			</s:Code>
			// 			<s:Reason>
			// 				<s:Text xml:lang="en-US">This User is not authorized to enroll</s:text>
			// 			</s:Reason>
			// 		</s:Fault>
			// 	</s:Body>
			// </s:Envelope>`)
			/* END TEMP */
			w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			w.Write(response)
		}
	}
}
