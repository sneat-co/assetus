import { Component, inject } from '@angular/core';
import {
  FormControl,
  UntypedFormGroup,
  Validators,
} from '@angular/forms';
import { ISpaceContext } from '@sneat/space-models';
import { SpaceBaseComponent } from '@sneat/space-components';
import {
  ASSET_SERVICE,
  AssetExtraType,
  IAssetDboBase,
  IAssetExtra,
  IAssetService,
  ICreateAssetRequest,
} from '@sneat/extension-assetus-contract';

// Lib-local base component for the three asset-add components. Ported from the
// legacy legacy ext-assetus-components AddAssetBaseComponent so the components
// no longer depend on the legacy lib. Members mirror the legacy base exactly;
// only the create flow is adapted to the lib's own AssetService (which hits the
// live assetus/create_asset backend with a FLAT ICreateAssetRequest, replacing
// the deleted legacy assets/create_asset envelope endpoint).
@Component({ template: '' })
export abstract class AddAssetBaseComponent extends SpaceBaseComponent {
  public static readonly metadata = { inputs: ['space'] };

  public isSubmitting = false;
  public titleForm = new UntypedFormGroup({
    title: new FormControl<string>('', Validators.required),
  });

  protected readonly assetService: IAssetService = inject(ASSET_SERVICE);

  // Accepts the legacy `{ asset, spaceID }` envelope the add-asset components
  // build, translates it into the lib's flat ICreateAssetRequest, creates the
  // asset via the lib AssetService, then navigates to the new asset's page.
  protected createAssetAndGoToAssetPage<
    ExtraType extends AssetExtraType,
    Extra extends IAssetExtra,
  >(
    request: {
      asset: IAssetDboBase<ExtraType, Extra>;
      spaceID: string;
      memberID?: string;
    },
    space: ISpaceContext,
  ): void {
    if (!this.space) {
      throw new Error('no team context');
    }

    const asset = request.asset;
    const flatRequest: ICreateAssetRequest = {
      spaceID: request.spaceID,
      name: asset.title || asset.name || '',
      category: asset.category,
      condition: asset.condition ?? 'good',
      status: asset.status,
      visibility: asset.visibility,
      type: asset.type,
      possession: asset.possession,
      yearOfBuild: asset.yearOfBuild,
      extraType: asset.extraType,
      extra: asset.extra as Record<string, unknown> | undefined,
      countryID: asset.countryID,
    };

    this.assetService.createAsset(flatRequest).subscribe({
      next: (resp) => {
        this.spaceParams.spaceNavService
          .navigateForwardToSpacePage(space, 'asset/' + resp.id, {
            replaceUrl: true,
            state: { asset: resp.asset, space },
          })
          .catch(
            this.spaceParams.errorLogger.logErrorHandler(
              'failed to navigate to team page',
            ),
          );
      },
      error: (err) => {
        this.isSubmitting = false;
        this.spaceParams.errorLogger.logError(
          err,
          'Failed to create a new asset',
        );
      },
    });
  }
}
