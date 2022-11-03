import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NewUpdateComponent } from './new-update.component';

describe('NewUpdateComponent', () => {
    let component: NewUpdateComponent;
    let fixture: ComponentFixture<NewUpdateComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            declarations: [ NewUpdateComponent ]
        })
            .compileComponents();

        fixture = TestBed.createComponent(NewUpdateComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
