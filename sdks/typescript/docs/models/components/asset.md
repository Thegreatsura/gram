# Asset

## Example Usage

```typescript
import { Asset } from "@gram/sdk/models/components";

let value: Asset = {
  contentLength: 95314,
  contentType: "<value>",
  createdAt: new Date("2025-11-12T07:39:23.206Z"),
  id: "<id>",
  kind: "openapiv3",
  sha256: "<value>",
  updatedAt: new Date("2023-04-08T11:11:06.209Z"),
  url: "https://wretched-humor.biz",
};
```

## Fields

| Field                                                                                         | Type                                                                                          | Required                                                                                      | Description                                                                                   |
| --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| `contentLength`                                                                               | *number*                                                                                      | :heavy_check_mark:                                                                            | The content length of the asset                                                               |
| `contentType`                                                                                 | *string*                                                                                      | :heavy_check_mark:                                                                            | The content type of the asset                                                                 |
| `createdAt`                                                                                   | [Date](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date) | :heavy_check_mark:                                                                            | The creation date of the asset.                                                               |
| `id`                                                                                          | *string*                                                                                      | :heavy_check_mark:                                                                            | The ID of the asset                                                                           |
| `kind`                                                                                        | [components.Kind](../../models/components/kind.md)                                            | :heavy_check_mark:                                                                            | N/A                                                                                           |
| `sha256`                                                                                      | *string*                                                                                      | :heavy_check_mark:                                                                            | The SHA256 hash of the asset                                                                  |
| `updatedAt`                                                                                   | [Date](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date) | :heavy_check_mark:                                                                            | The last update date of the asset.                                                            |
| `url`                                                                                         | *string*                                                                                      | :heavy_check_mark:                                                                            | The URL to the uploaded asset                                                                 |