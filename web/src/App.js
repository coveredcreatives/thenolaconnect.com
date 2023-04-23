import * as React from 'react';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Outlet,
} from "react-router-dom";
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Layout } from 'antd';
import { Navigation } from './components/navigation';
import { Footer } from './components/footer';
import Pages from './pages'

const queryClient = new QueryClient()

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <Routes>
          <Route path="/" element={<FullLayout />}>
            <Route index element={<Pages.Welcome />}></Route>
            <Route path="/qr-code" element={<Pages.GenerateQRCode />}></Route>
            <Route path="/orders" element={<Pages.OrderReview />}></Route>
            <Route path="/orders/placement" element={<Pages.OrderPlacement />}></Route>
            <Route path="/tos" element={<Pages.Welcome />}></Route>
            <Route path="/privacy" element={<Pages.Welcome />}></Route>
          </Route>
        </Routes>
      </Router>
    </QueryClientProvider>
  );
}

function FullLayout() {
  return (
    <div className="App">
      <Layout.Header>
        <Navigation />
      </Layout.Header>
      <Outlet />
      <Layout.Footer>
        <Footer />
      </Layout.Footer>
    </div>
  )
}

export default App;
