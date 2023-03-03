import * as React from 'react';
import {
  Box,
  Heading,
  PageLayout,
} from '@primer/react'


export function Welcome() {
    return (
        <PageLayout>
          <PageLayout.Content>
            <Box display="grid" gridGap={3}>
              <Heading>Pardon Our Progress</Heading>
              <p>This is the future home to a covered creatives administration tool.</p>
            </Box>
          </PageLayout.Content>
        </PageLayout>
    )
  }