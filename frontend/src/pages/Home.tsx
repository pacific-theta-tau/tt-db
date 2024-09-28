// Home.tsx (default export)
import React from 'react';
import { Link } from 'react-router-dom';

const HomePage: React.FC = () => {
  return (
      <div>
          <h1>Home Page</h1>
          <Link to="/brothers">See all members</Link>
          <Link to="/events">See all events</Link>
      </div>
  );
};

export default HomePage;
