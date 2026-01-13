import React from 'react';
import Header from '../components/Header';
import Footer from '../components/Footer';
import HeroSection from '../components/HeroSection.tsx';
import ServicesSection from '../components/ServicesSection.tsx';
import StatsSection from '../components/StatsSection.tsx';

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-grow">
        <HeroSection />
        <ServicesSection />
        <StatsSection />
      </main>
      <Footer />
    </div>
  );
}