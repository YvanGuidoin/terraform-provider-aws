// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package ecs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// GetTag fetches an individual ecs service tag for a resource.
// Returns whether the key value and any errors. A NotFoundError is used to signal that no value was found.
// This function will optimise the handling over ListTags, if possible.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func GetTag(conn ecsiface.ECSAPI, identifier string, key string) (*string, error) {
	return GetTagWithContext(context.Background(), conn, identifier, key)
}
func GetTagWithContext(ctx context.Context, conn ecsiface.ECSAPI, identifier string, key string) (*string, error) {
	listTags, err := ListTagsWithContext(ctx, conn, identifier)

	if err != nil {
		return nil, err
	}

	if !listTags.KeyExists(key) {
		return nil, tfresource.NewEmptyResultError(nil)
	}

	return listTags.KeyValue(key), nil
}

// ListTags lists ecs service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(conn ecsiface.ECSAPI, identifier string) (tftags.KeyValueTags, error) {
	return ListTagsWithContext(context.Background(), conn, identifier)
}

func ListTagsWithContext(ctx context.Context, conn ecsiface.ECSAPI, identifier string) (tftags.KeyValueTags, error) {
	input := &ecs.ListTagsForResourceInput{
		ResourceArn: aws.String(identifier),
	}

	output, err := conn.ListTagsForResourceWithContext(ctx, input)

	if tfawserr.ErrMessageContains(err, "InvalidParameterException", "The specified cluster is inactive. Specify an active cluster and try again.") {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.Tags), nil
}

// []*SERVICE.Tag handling

// Tags returns ecs service tags.
func Tags(tags tftags.KeyValueTags) []*ecs.Tag {
	result := make([]*ecs.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &ecs.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from ecs service tags.
func KeyValueTags(tags []*ecs.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}

// UpdateTags updates ecs service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn ecsiface.ECSAPI, identifier string, oldTags interface{}, newTags interface{}) error {
	return UpdateTagsWithContext(context.Background(), conn, identifier, oldTags, newTags)
}
func UpdateTagsWithContext(ctx context.Context, conn ecsiface.ECSAPI, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &ecs.UntagResourceInput{
			ResourceArn: aws.String(identifier),
			TagKeys:     aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.UntagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &ecs.TagResourceInput{
			ResourceArn: aws.String(identifier),
			Tags:        Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.TagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
