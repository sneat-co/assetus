import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { ModalController } from '@ionic/angular/standalone';
import { AssetsPageComponent } from './assets-page.component';
import { assetusTestProviders } from '../../test-providers';

describe('AssetsPageComponent', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [AssetsPageComponent],
      providers: [
        provideRouter([]),
        ...assetusTestProviders(),
        { provide: ModalController, useValue: { create: () => Promise.resolve({}) } },
      ],
    }),
  );

  it('creates and renders the assets toolbar', () => {
    const fixture = TestBed.createComponent(AssetsPageComponent);
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    expect(host.querySelector('ion-title')?.textContent).toContain('Assets');
  });
});
