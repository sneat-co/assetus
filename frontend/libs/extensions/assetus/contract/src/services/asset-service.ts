import { InjectionToken } from '@angular/core';
import { ISpaceContext } from '@sneat/space-models';
import { Observable } from 'rxjs';
import { IAssetContext } from '../contexts';
import { IAssetDbo } from '../dto';
import {
  IAssetHistoryResponse,
  IAssetResponse,
  ICreateAssetRequest,
  ICreateVehicleRecordRequest,
  ICreateVehicleRecordResponse,
  IRecordHistoryEventRequest,
  IRemoveAssetRequest,
  ITransferAssetRequest,
  ITransferAssetResponse,
  IUpdateAssetRequest,
} from './interfaces';

// An asset document plus its Firestore id (the `id` field is merged in by
// collectionData's idField option). Pure data shape consumed by watchAssets().
export interface IIdAndAssetDbo {
  id: string;
  dbo: IAssetDbo;
}

// IAssetService is the runtime-light contract the asset components and pages
// depend on. Members mirror the concrete AssetService's public surface exactly;
// the implementation lives in the internal/shared lib (moved in a later task)
// and is provided via the ASSET_SERVICE token below.
export interface IAssetService {
  createAsset(request: ICreateAssetRequest): Observable<IAssetResponse>;
  updateAsset(request: IUpdateAssetRequest): Observable<IAssetResponse>;
  removeAsset(request: IRemoveAssetRequest): Observable<void>;
  transferAsset(
    request: ITransferAssetRequest,
  ): Observable<ITransferAssetResponse>;
  recordHistoryEvent(request: IRecordHistoryEventRequest): Observable<void>;
  getHistory(
    spaceID: string,
    assetID: string,
  ): Observable<IAssetHistoryResponse>;
  addVehicleRecord(
    request: ICreateVehicleRecordRequest,
  ): Observable<ICreateVehicleRecordResponse>;
  watchAssets(spaceID: string): Observable<IIdAndAssetDbo[]>;
  watchAssetByID(
    space: ISpaceContext,
    id: string,
  ): Observable<IAssetContext>;
}

export const ASSET_SERVICE = new InjectionToken<IAssetService>('AssetService');
