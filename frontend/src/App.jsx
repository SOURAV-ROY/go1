import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, useSearchParams, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { CreditCard, CheckCircle, XCircle, Loader2, Moon, Sun } from 'lucide-react';

const API_BASE = 'http://localhost:8080/api';

const ThemeToggle = () => {
  const [isDark, setIsDark] = useState(true);

  useEffect(() => {
    if (isDark) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [isDark]);

  return (
    <button
      onClick={() => setIsDark(!isDark)}
      className="p-2 rounded-lg bg-zinc-200 dark:bg-zinc-800 text-zinc-800 dark:text-zinc-200 transition-colors"
    >
      {isDark ? <Sun size={20} /> : <Moon size={20} />}
    </button>
  );
};

const CheckoutCard = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handlePayment = async () => {
    setLoading(true);
    setError(null);
    try {
      const resp = await axios.post(`${API_BASE}/pay`, {
        orderId: `ORD-${Math.floor(Math.random() * 1000000)}`,
        amount: "100.00",
        ip: "127.0.0.1"
      });
      
      if (resp.data.callBackUrl) {
        window.location.href = resp.data.callBackUrl;
      } else {
        setError("Failed to get redirection URL");
      }
    } catch (err) {
      setError(err.response?.data || err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-md w-full p-8 bg-white dark:bg-[#0c0c0f] border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-xl transition-all">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold tracking-tight text-zinc-950 dark:text-zinc-50">Checkout</h2>
        <CreditCard className="text-zinc-400" />
      </div>
      
      <div className="space-y-4 mb-8">
        <div className="flex justify-between text-zinc-500 dark:text-zinc-400">
          <span>Product</span>
          <span className="text-zinc-950 dark:text-zinc-50 font-medium">Premium Subscription</span>
        </div>
        <div className="flex justify-between text-zinc-500 dark:text-zinc-400">
          <span>Amount</span>
          <span className="text-zinc-950 dark:text-zinc-50 font-medium">100.00 BDT</span>
        </div>
        <div className="h-px bg-zinc-100 dark:bg-zinc-800 w-full" />
        <div className="flex justify-between text-lg font-bold">
          <span>Total</span>
          <span className="text-emerald-600">100.00 BDT</span>
        </div>
      </div>

      {error && (
        <div className="mb-4 p-3 bg-rose-50 dark:bg-rose-900/20 border border-rose-200 dark:border-rose-800/30 text-rose-800 dark:text-rose-300 rounded-lg text-sm">
          {error}
        </div>
      )}

      <button
        onClick={handlePayment}
        disabled={loading}
        className="w-full bg-[#059669] hover:bg-[#047857] disabled:opacity-50 disabled:cursor-not-allowed text-white font-semibold py-3 rounded-xl transition-all shadow-lg flex items-center justify-center gap-2"
      >
        {loading ? <Loader2 className="animate-spin" size={20} /> : null}
        {loading ? 'Initializing...' : 'Pay with Nagad'}
      </button>
      
      <p className="mt-4 text-center text-xs text-zinc-500 dark:text-zinc-500">
        You will be redirected to Nagad's secure payment portal
      </p>
    </div>
  );
};

const PaymentCallback = () => {
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState('verifying'); // verifying, success, failed
  const navigate = useNavigate();

  useEffect(() => {
    const verifyPayment = async () => {
      const paymentRefId = searchParams.get('paymentRefId');
      const apiStatus = searchParams.get('status');

      if (apiStatus === 'Success' && paymentRefId) {
        try {
          const resp = await axios.post(`${API_BASE}/verify`, { paymentRefId });
          if (resp.data.status === 'Success') {
            setStatus('success');
          } else {
            setStatus('failed');
          }
        } catch (err) {
          setStatus('failed');
        }
      } else {
        setStatus('failed');
      }
    };

    verifyPayment();
  }, [searchParams]);

  return (
    <div className="max-w-md w-full p-8 bg-white dark:bg-[#0c0c0f] border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-xl text-center">
      {status === 'verifying' && (
        <div className="space-y-4">
          <Loader2 className="animate-spin mx-auto text-emerald-500" size={48} />
          <h2 className="text-xl font-bold">Verifying Payment...</h2>
          <p className="text-zinc-500">Please wait while we confirm your transaction with Nagad.</p>
        </div>
      )}

      {status === 'success' && (
        <div className="space-y-4">
          <CheckCircle className="mx-auto text-emerald-500" size={48} />
          <h2 className="text-xl font-bold text-zinc-950 dark:text-zinc-50">Payment Successful!</h2>
          <p className="text-zinc-500">Your order has been processed. Thank you for your purchase.</p>
          <button 
            onClick={() => navigate('/')}
            className="mt-6 text-emerald-600 hover:text-emerald-700 font-medium"
          >
            Go back to Home
          </button>
        </div>
      )}

      {status === 'failed' && (
        <div className="space-y-4">
          <XCircle className="mx-auto text-rose-500" size={48} />
          <h2 className="text-xl font-bold text-zinc-950 dark:text-zinc-50">Payment Failed</h2>
          <p className="text-zinc-500">Something went wrong or the payment was cancelled.</p>
          <button 
            onClick={() => navigate('/')}
            className="mt-6 text-rose-600 hover:text-rose-700 font-medium"
          >
            Try Again
          </button>
        </div>
      )}
    </div>
  );
};

const App = () => {
  return (
    <Router>
      <div className="min-h-screen flex flex-col items-center justify-center p-6 transition-colors duration-300">
        <div className="fixed top-6 right-6">
          <ThemeToggle />
        </div>
        
        <Routes>
          <Route path="/" element={<CheckoutCard />} />
          <Route path="/callback" element={<PaymentCallback />} />
        </Routes>
      </div>
    </Router>
  );
};

export default App;
