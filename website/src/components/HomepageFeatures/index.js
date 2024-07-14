import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'Optimized Router',
    Svg: require('@site/static/img/0568-rocket.svg').default,
    description: (
      <>
        Echo boasts a highly optimized HTTP router that operates without dynamic memory allocation.
        This router intelligently prioritizes routes for efficient routing and processing of incoming HTTP requests.
        The absence of dynamic memory allocation contributes to better performance and resource utilization.
      </>
    ),
  },
  {
    title: 'Scalable',
    Svg: require('@site/static/img/0075-cloud.svg').default,
    description: (
      <>
        Echo framework empowers developers to build robust and scalable RESTful APIs effortlessly.
        It facilitates the organization of endpoints into logical groups, simplifying the management of complex APIs.
        This scalability ensures that your application can grow and handle increased traffic or functionality seamlessly.
      </>
    ),
  },
  {
    title: 'Automatic TLS',
    Svg: require('@site/static/img/0221-license2.svg').default,
    description: (
      <>
        Echo streamlines the process of enabling secure communication by automatically handling TLS certificate installation from Let's Encrypt.
        This automation simplifies the setup of secure connections, allowing you to secure your application with minimal effort.
      </>
    ),
  },
  {
    title: 'HTTP/2 Support',
    Svg: require('@site/static/img/0567-speed-fast.svg').default,
    description: (
      <>
        Echo embraces the HTTP/2 protocol, which enhances the speed and responsiveness of your web applications.
        HTTP/2 is designed to provide a faster and more efficient user experience by optimizing the way data is transmitted between the server and the client.
      </>
    ),
  },
  {
    title: 'Middleware',
    Svg: require('@site/static/img/0893-funnel.svg').default,
    description: (
      <>
        Echo offers a rich collection of built-in middleware functions that can be utilized to enhance your application's functionality and security.
        You can apply middleware at different levels: globally at the root level, within specific route groups, or on individual routes.
        Additionally, Echo allows developers to define their custom middleware to tailor the application's behavior to specific needs.
      </>
    ),
  },
  {
    title: 'Data Binding',
    Svg: require('@site/static/img/0101-database-upload.svg').default,
    description: (
      <>
        Echo simplifies the process of binding HTTP request payloads, including JSON, XML, or form-data.
        This feature makes it easy to extract and work with data sent in the request body, enabling seamless integration with your application logic.
      </>
    ),
  },
  {
    title: 'Data Rendering',
    Svg: require('@site/static/img/0102-database-download.svg').default,
    description: (
      <>
        Echo provides a versatile API for sending various types of HTTP responses.
        You can easily send responses in formats such as JSON, XML, HTML, files, attachments, inline content, streams, or blobs.
        This flexibility ensures that your application can generate and serve different types of content to clients based on their needs.
      </>
    ),
  },
  {
    title: 'Templates',
    Svg: require('@site/static/img/0780-code.svg').default,
    description: (
      <>
        Echo supports template rendering using any template engine of your choice.
        This feature allows you to generate dynamic HTML content for rendering web pages, email templates, or other view-based responses.
        You have the freedom to select the template engine that best suits your project.
      </>
    ),
  },
  {
    title: 'Extensible',
    Svg: require('@site/static/img/0243-equalizer.svg').default,
    description: (
      <>
        Echo offers extensibility at various levels.
        You can customize centralized HTTP error handling to gracefully manage errors and exceptions.
        Moreover, the framework provides an easily extendable API, allowing you to create custom middleware, plugins,
        or components to tailor Echo to your specific requirements and integrate third-party functionality seamlessly.
      </>
    ),
  }
];

function Feature({ Svg, title, description }) {
  return (
    <div className={clsx('col col--4')}>
      <div className={clsx('card', 'text--left', 'padding-horiz--md', styles.featureCard)}>
        <div className={clsx('flex', 'flex--center', 'margin-bottom--sm','padding-top--md', styles.headline)}>
          <Svg className={styles.featureSvg} role="img" />
          <h3 className="margin--none">{title}</h3>
        </div>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <h1 className="text--left">Features</h1>
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
