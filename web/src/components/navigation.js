import * as React from 'react';
import {
    UnderlineNav,
} from '@primer/react';
import { NavLink } from 'react-router-dom';

export function Navigation() {
    return (
        <UnderlineNav aria-label="Main">
            <UnderlineNav.Link to="/" as={NavLink} className={(navData) => (navData.isActive ? 'active' : 'none')}>
                Home
            </UnderlineNav.Link>
            <UnderlineNav.Link to="/qr-code" as={NavLink} className={(navData) => (navData.isActive ? 'active' : 'none')}>
                QR Management
            </UnderlineNav.Link>
            <UnderlineNav.Link to="/orders" as={NavLink} className={(navData) => (navData.isActive ? 'active' : 'none')}>
                Orders
            </UnderlineNav.Link>
            <UnderlineNav.Link href="https://docs.google.com/forms/d/e/1FAIpQLSeKe8iSippG-8wLxdPaGrL2Bpbqw6O8lofNN6gti98MX-YzOw/viewform" as="a">
                Order Placement
            </UnderlineNav.Link>
        </UnderlineNav>
    )
}