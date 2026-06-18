import { ISpaceItemNavContext } from '@sneat/space-models';
import { IAssetBrief, IAssetDbo, IOwner } from '../dto';

export interface IAssetContext
  extends ISpaceItemNavContext<IAssetBrief, IAssetDbo> {
  owner?: IOwner;
}
