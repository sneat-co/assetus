import { HttpParams } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import {
  Firestore as AngularFirestore,
  CollectionReference,
  collection,
  collectionData,
} from '@angular/fire/firestore';
import { SneatApiService } from '@sneat/api';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { IAssetDbo } from '../dto';
import {
  ICreateVehicleRecordRequest,
  ICreateVehicleRecordResponse,
  IAssetHistoryResponse,
  IAssetResponse,
  ICreateAssetRequest,
  IGetAssetResponse,
  IRecordHistoryEventRequest,
  IRemoveAssetRequest,
  ITransferAssetRequest,
  ITransferAssetResponse,
  IUpdateAssetRequest,
} from './interfaces';

// An asset document plus its Firestore id (the `id` field is merged in by
// collectionData's idField option). Local alias to avoid guessing a generic
// @sneat/core helper name.
export interface IIdAndAssetDbo {
  id: string;
  dbo: IAssetDbo;
}

// All assetus backend endpoints live under the `assetus/` path on the shared
// sneat API (the SneatApiService prefixes the v0 base URL).
const api = (endpoint: string): string => `assetus/${endpoint}`;

@Injectable()
export class AssetService {
  private readonly afs = inject(AngularFirestore);
  private readonly sneatApiService = inject(SneatApiService);

  public createAsset(
    request: ICreateAssetRequest,
  ): Observable<IAssetResponse> {
    return this.sneatApiService.post<IAssetResponse>(
      api('create_asset'),
      request,
    );
  }

  public getAsset(
    spaceID: string,
    assetID: string,
  ): Observable<IGetAssetResponse> {
    const params = new HttpParams({ fromObject: { spaceID, assetID } });
    return this.sneatApiService.get<IGetAssetResponse>(api('asset'), params);
  }

  public updateAsset(
    request: IUpdateAssetRequest,
  ): Observable<IAssetResponse> {
    return this.sneatApiService.post<IAssetResponse>(
      api('update_asset'),
      request,
    );
  }

  public removeAsset(request: IRemoveAssetRequest): Observable<void> {
    return this.sneatApiService.post<void>(api('remove_asset'), request);
  }

  public transferAsset(
    request: ITransferAssetRequest,
  ): Observable<ITransferAssetResponse> {
    return this.sneatApiService.post<ITransferAssetResponse>(
      api('transfer_asset'),
      request,
    );
  }

  public recordHistoryEvent(
    request: IRecordHistoryEventRequest,
  ): Observable<void> {
    return this.sneatApiService.post<void>(
      api('record_history_event'),
      request,
    );
  }

  public getHistory(
    spaceID: string,
    assetID: string,
  ): Observable<IAssetHistoryResponse> {
    const params = new HttpParams({ fromObject: { spaceID, assetID } });
    return this.sneatApiService.get<IAssetHistoryResponse>(
      api('asset_history'),
      params,
    );
  }

  public addVehicleRecord(
    request: ICreateVehicleRecordRequest,
  ): Observable<ICreateVehicleRecordResponse> {
    return this.sneatApiService.post<ICreateVehicleRecordResponse>(
      api('create_vehicle_record'),
      request,
    );
  }

  // Live read of a space's assets straight from Firestore (reads the collection
  // directly rather than via HTTP). The backend writes assets to
  // spaces/{spaceID}/ext/assetus/assets/{assetID}.
  public watchAssets(spaceID: string): Observable<IIdAndAssetDbo[]> {
    return collectionData<IAssetDbo>(this.assetsCollection(spaceID), {
      idField: 'id',
    }).pipe(
      map((dbos) => dbos.map((dbo) => ({ id: dbo.id, dbo }))),
    );
  }

  private assetsCollection(spaceID: string): CollectionReference<IAssetDbo> {
    return collection(
      this.afs,
      'spaces',
      spaceID,
      'ext',
      'assetus',
      'assets',
    ) as CollectionReference<IAssetDbo>;
  }
}
