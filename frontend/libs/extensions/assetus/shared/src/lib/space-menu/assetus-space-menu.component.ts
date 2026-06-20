import { TitleCasePipe } from '@angular/common';
import {
  ChangeDetectionStrategy,
  Component,
  computed,
  inject,
  signal,
} from '@angular/core';
import { RouterLink } from '@angular/router';
import {
  IonIcon,
  IonItem,
  IonLabel,
  IonList,
  IonSelect,
  IonSelectOption,
  MenuController,
} from '@ionic/angular/standalone';
import { ISneatUserState } from '@sneat/auth-core';
import { IUserSpaceBrief } from '@sneat/auth-models';
import { AuthMenuItemComponent } from '@sneat/auth-ui';
import { IIdAndBrief } from '@sneat/core';
import {
  SpaceBaseComponent,
  SpaceComponentBaseParams,
} from '@sneat/space-components';
import { SpaceServiceModule } from '@sneat/space-services';
import { zipMapBriefsWithIDs } from '@sneat/space-models';
import { ClassName } from '@sneat/ui';
import { takeUntil } from 'rxjs/operators';

// assetus-specific side menu: a space selector (to switch spaces) plus the
// Assets entry for the active space. A focused space menu for the
// ownership-only asset domain.
@Component({
  selector: 'assetus-space-menu',
  templateUrl: './assetus-space-menu.component.html',
  imports: [
    TitleCasePipe,
    RouterLink,
    SpaceServiceModule,
    IonList,
    IonItem,
    IonSelect,
    IonSelectOption,
    IonIcon,
    IonLabel,
    AuthMenuItemComponent,
  ],
  providers: [
    { provide: ClassName, useValue: 'AssetusSpaceMenuComponent' },
    SpaceComponentBaseParams,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AssetusSpaceMenuComponent extends SpaceBaseComponent {
  protected readonly $spaces = signal<
    readonly IIdAndBrief<IUserSpaceBrief>[] | undefined
  >(undefined);
  protected readonly $disabled = computed(() => !this.$spaceID());

  private readonly menuCtrl = inject(MenuController);

  constructor() {
    super();
    this.spaceParams.userService.userState
      .pipe(takeUntil(this.destroyed$))
      .subscribe({
        next: (userState: ISneatUserState) =>
          this.$spaces.set(
            userState?.record
              ? zipMapBriefsWithIDs(userState.record.spaces) || []
              : undefined,
          ),
        error: this.errorLogger.logErrorHandler('failed to get user state'),
      });
  }

  protected onSpaceSelected(event: Event): void {
    const spaceID = (event as CustomEvent).detail.value as string;
    if (spaceID === this.space?.id) {
      return;
    }
    const space = this.$spaces()?.find((t) => t.id === spaceID);
    if (space) {
      this.setSpaceRef(space);
      this.spaceNav
        .navigateToSpace(space)
        .catch(
          this.errorLogger.logErrorHandler(
            'Failed to navigate to selected space',
          ),
        );
    }
    this.closeMenu();
  }

  protected closeMenu(): void {
    this.menuCtrl.close().catch(this.errorLogger.logError);
  }
}
