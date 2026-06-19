import { Component } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { describe, expect, it, vi } from 'vitest';
import { of } from 'rxjs';
import { ClassName } from '@sneat/ui';
import { SpaceComponentBaseParams } from '@sneat/space-components';
import { ISpaceContext } from '@sneat/space-models';
import { spacePageTestProviders } from '../../../testing/test-providers';
import { AssetService } from '../../services';
import { AddAssetBaseComponent } from './add-asset-base.component';

// Trivial concrete subclass exposing the protected create flow so the abstract
// base can be exercised through TestBed.
@Component({
  template: '',
  providers: [
    { provide: ClassName, useValue: 'TestAddAssetComponent' },
    SpaceComponentBaseParams,
  ],
})
class TestAddAssetComponent extends AddAssetBaseComponent {
  public callCreate(
    request: Parameters<TestAddAssetComponent['createAssetAndGoToAssetPage']>[0],
    space: ISpaceContext,
  ): void {
    this.createAssetAndGoToAssetPage(request, space);
  }
}

describe('AddAssetBaseComponent', () => {
  const setup = () => {
    const createAsset = vi
      .fn()
      .mockReturnValue(of({ id: 'new1', asset: { id: 'new1' } }));
    const navigate = vi.fn(
      (..._args: unknown[]) => Promise.resolve(true),
    );

    TestBed.configureTestingModule({
      imports: [TestAddAssetComponent],
      providers: [
        ...spacePageTestProviders(),
        { provide: AssetService, useValue: { createAsset } },
      ],
    });
    const fixture = TestBed.createComponent(TestAddAssetComponent);
    const component = fixture.componentInstance;
    // Seed the read-only `space` getter via the protected $spaceRef signal,
    // and stub the nav service the create flow calls on success.
    (
      component as unknown as { $spaceRef: { set(v: unknown): void } }
    ).$spaceRef.set({ id: 's1' });
    (
      component as unknown as {
        spaceParams: { spaceNavService: { navigateForwardToSpacePage: unknown } };
      }
    ).spaceParams.spaceNavService.navigateForwardToSpacePage = navigate;

    return { component, createAsset, navigate };
  };

  it('translates the {asset,spaceID} envelope into a flat ICreateAssetRequest', () => {
    const { component, createAsset } = setup();
    const space = { id: 's1' } as ISpaceContext;

    component.callCreate(
      {
        spaceID: 's1',
        asset: {
          title: 'My car',
          category: 'vehicle',
          extraType: 'vehicle',
          extra: { make: 'Toyota', model: 'Corolla' },
          yearOfBuild: 2020,
        } as never,
      },
      space,
    );

    expect(createAsset).toHaveBeenCalledOnce();
    const flat = createAsset.mock.calls[0][0];
    expect(flat).toMatchObject({
      spaceID: 's1',
      name: 'My car',
      category: 'vehicle',
      condition: 'good',
      extraType: 'vehicle',
      extra: { make: 'Toyota', model: 'Corolla' },
      yearOfBuild: 2020,
    });
  });

  it('navigates to the new asset page on success', () => {
    const { component, navigate } = setup();
    const space = { id: 's1' } as ISpaceContext;

    component.callCreate(
      { spaceID: 's1', asset: { category: 'vehicle' } as never },
      space,
    );

    expect(navigate).toHaveBeenCalledOnce();
    expect(navigate.mock.calls[0][0]).toBe(space);
    expect(navigate.mock.calls[0][1]).toBe('asset/new1');
  });
});
