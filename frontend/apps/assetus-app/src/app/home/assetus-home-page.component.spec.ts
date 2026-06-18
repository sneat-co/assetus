import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { getStandardSneatProviders } from '@sneat/app';
import { SneatUserService } from '@sneat/auth-core';
import { BehaviorSubject } from 'rxjs';
import { assetusAppEnvironmentConfig } from '../../environments/environment';
import { AssetusHomePageComponent } from './assetus-home-page.component';

// Renders the home page for an AUTHENTICATED user who HAS spaces, so the
// embedded SpacesCard -> SpacesList chain is actually constructed. That chain
// needs SpaceService + UserRequiredFieldsService, which only surface at runtime
// as NG0201. A signed-out user would not render the list and would miss this.
describe('AssetusHomePageComponent', () => {
  const userState$ = new BehaviorSubject<unknown>({
    status: 'authenticated',
    user: { uid: 'u1', isAnonymous: false, emailVerified: true, providerData: [] },
    record: {
      title: 'Test User',
      spaces: { s1: { title: 'Family', type: 'family', roles: ['creator'] } },
    },
  });

  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [AssetusHomePageComponent],
      providers: [
        ...getStandardSneatProviders(assetusAppEnvironmentConfig),
        provideRouter([]),
        {
          provide: SneatUserService,
          useValue: { userState: userState$, currentUserID: 'u1' },
        },
      ],
    }),
  );

  it('renders the spaces list for a user with spaces (all DI resolves, no NG0201)', () => {
    const fixture = TestBed.createComponent(AssetusHomePageComponent);
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    expect(host.querySelector('sneat-spaces-card')).toBeTruthy();
    expect(host.querySelector('sneat-spaces-list')).toBeTruthy();
  });
});
