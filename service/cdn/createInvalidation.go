package cdn

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/gofiber/fiber/v2"
	"time"
)

type CreateInvalidationBody struct {
	DomainName string `json:"domainName"`
	UrlPath    string `json:"urlPath"`
}

func CreateInvalidation(c *fiber.Ctx) error {
	var distributionID string
	createInvalidationBody := CreateInvalidationBody{}
	_ = json.Unmarshal(c.Body(), &createInvalidationBody)

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2")},
	)

	cf := cloudfront.New(sess)
	result, err := cf.ListDistributions(nil)
	if err != nil {
		return c.JSON(fiber.Map{"message": "cf ListDsistribution occured Error", "error": err})
	}

	for _, distribution := range result.DistributionList.Items {
		if distribution.Aliases != nil {
			for _, alias := range distribution.Aliases.Items {
				if *alias == createInvalidationBody.DomainName {
					distributionID = *distribution.Id
					break
				}
			}
		}
	}
	fmt.Println(distributionID)

	if distributionID == "" {
		err := fmt.Errorf("cannot found cloudfront distributionID")
		return c.JSON(err)
	}

	invalidationInput := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distributionID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			Paths: &cloudfront.Paths{
				Items:    []*string{aws.String(createInvalidationBody.UrlPath)},
				Quantity: aws.Int64(1),
			},
			CallerReference: aws.String(time.Now().String()),
		},
	}

	invalRe, invalErr := cf.CreateInvalidation(invalidationInput)
	if invalErr != nil {
		return c.JSON(invalErr)
	}
	fmt.Println("Invalidation created:", *invalRe.Invalidation.Id)

	return c.JSON(fiber.Map{"invalidation craeted": *invalRe.Invalidation.Id})
}
