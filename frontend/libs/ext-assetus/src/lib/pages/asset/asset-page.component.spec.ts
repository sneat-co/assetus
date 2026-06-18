import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import {
  AlertController,
  ModalController,
} from '@ionic/angular/standalone';
import { AssetPageComponent } from './asset-page.component';
import { assetusTestProviders } from '../../test-providers';

describe('AssetPageComponent', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [AssetPageComponent],
      providers: [
        provideRouter([]),
        ...assetusTestProviders(),
        { provide: ModalController, useValue: { create: () => Promise.resolve({}) } },
        { provide: AlertController, useValue: { create: () => Promise.resolve({}) } },
      ],
    }),
  );

  it('creates and renders the edit form', () => {
    const fixture = TestBed.createComponent(AssetPageComponent);
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    expect(host.querySelector('ion-content')).toBeTruthy();
    expect(host.querySelector('assetus-asset-history-timeline')).toBeTruthy();
  });
});
