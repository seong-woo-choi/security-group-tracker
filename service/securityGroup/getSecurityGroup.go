package securityGroup

import (
	alb "go-sdk/service/securityGroup/alb"
	docdb "go-sdk/service/securityGroup/docdb"
	ec2 "go-sdk/service/securityGroup/ec2"
	elasticache "go-sdk/service/securityGroup/elasticache"
	msk "go-sdk/service/securityGroup/msk"
	rds "go-sdk/service/securityGroup/rds"

	"github.com/gofiber/fiber/v2"
)

func GetSecurityGroup(c *fiber.Ctx) error {
	resourceType := c.Query("type", "")
	resourceName := c.Query("name", "")

	switch resourceType {
	case "alb":
		err, alb := alb.GetSecurityGroup(resourceName)
		if err != nil {
			return c.JSON(fiber.Map{"message": err.Error()})
		}
		return c.JSON(alb)
	case "docdb":
		err, docdb := docdb.GetSecurityGroup(resourceName)
		if err != nil {
			return c.JSON(fiber.Map{"message": err.Error()})
		}
		return c.JSON(docdb)
	case "msk":
		err, msk := msk.GetSecurityGroup(resourceName)
		if err != nil {
			return c.JSON(fiber.Map{"message": err.Error()})
		}
		return c.JSON(msk)
	case "rds":
		err, rds := rds.GetSecurityGroup(resourceName)
		if err != nil {
			return c.JSON(fiber.Map{"message": err.Error()})
		}
		return c.JSON(rds)
	case "elasticache":
		err, elasticache := elasticache.GetSecurityGroup(resourceName)
		if err != nil {
			return c.JSON(fiber.Map{"message": err.Error()})
		}
		return c.JSON(elasticache)
	case "ec2":
		err, sg := ec2.GetSecurityGroup(resourceName)
		if err != nil {
			return c.JSON(fiber.Map{"message": err.Error()})
		}
		return c.JSON(sg)
	default:
		return c.JSON(fiber.Map{"status": "404", "message": "옳바른 리소스 타입을 입력해주세요(alb, elasticache, msk, rds, docdb, rds, sg)"})
	}
}
