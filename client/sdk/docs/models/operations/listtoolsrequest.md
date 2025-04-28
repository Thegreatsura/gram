# ListToolsRequest

## Example Usage

```typescript
import { ListToolsRequest } from "@gram/client/models/operations";

let value: ListToolsRequest = {};
```

## Fields

| Field                                                        | Type                                                         | Required                                                     | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `cursor`                                                     | *string*                                                     | :heavy_minus_sign:                                           | The cursor to fetch results from                             |
| `deploymentId`                                               | *string*                                                     | :heavy_minus_sign:                                           | The deployment ID. If unset, latest deployment will be used. |
| `gramSession`                                                | *string*                                                     | :heavy_minus_sign:                                           | Session header                                               |
| `gramProject`                                                | *string*                                                     | :heavy_minus_sign:                                           | project header                                               |