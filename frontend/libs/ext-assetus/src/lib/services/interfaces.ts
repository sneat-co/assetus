import {
  AssetCategory,
  AssetCondition,
  AssetVisibility,
  HistoryEventType,
  IAddVehicleRecordRequest as IAddVehicleRecord,
  IAssetDbo,
  IHistoryEvent,
  IMoney,
  IOwner,
} from '../dto';

// Optional metadata accepted by both create and update.
export interface IAssetMetadata {
  description?: string;
  acquisitionDate?: string;
  purchasePrice?: IMoney;
  estimatedValue?: IMoney;
  location?: string;
  notes?: string;
  tags?: string[];
}

export interface ICreateAssetRequest extends IAssetMetadata {
  spaceID: string;
  name: string;
  category: AssetCategory;
  condition: AssetCondition;
  // Optional override; when omitted the backend inherits the space default.
  visibility?: AssetVisibility;
}

export interface IAssetResponse {
  id: string;
  asset: IAssetDbo;
}

export interface IGetAssetResponse extends IAssetResponse {
  owner: IOwner;
}

export interface IUpdateAssetRequest extends IAssetMetadata {
  spaceID: string;
  assetID: string;
  name: string;
  category: AssetCategory;
  condition: AssetCondition;
  visibility: AssetVisibility;
}

export interface IRemoveAssetRequest {
  spaceID: string;
  assetID: string;
  // Soft-archive by default; hardDelete removes the record entirely.
  hardDelete?: boolean;
}

export interface ITransferAssetRequest {
  spaceID: string;
  assetID: string;
  toSpaceID: string;
}

export interface ITransferAssetResponse {
  id: string;
  owner: IOwner;
}

export interface IRecordHistoryEventRequest {
  spaceID: string;
  assetID: string;
  type: HistoryEventType;
  occurredAt?: string;
  note?: string;
}

export interface IAssetHistoryResponse {
  assetID: string;
  events: IHistoryEvent[];
}

// Appends a vehicle record (mileage and/or fuel reading) to a vehicle asset.
// The flat record shape (IAddVehicleRecord, from the dto) plus the spaceID
// required by the backend's embedded SpaceRequest.
export interface ICreateVehicleRecordRequest extends IAddVehicleRecord {
  spaceID: string;
}

export interface ICreateVehicleRecordResponse {
  id: string;
}
