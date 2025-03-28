package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/internal/server/proto"
	"strings"

	"github.com/olivere/elastic/v7"
)

type ElasticSearch struct {
	host string
	port int
}

func NewElasticService(host string, port int) *ElasticSearch {
	return &ElasticSearch{
		host: host,
		port: port,
	}
}

type FAQ struct {
    Question string `json:"question"`
    Answer   string `json:"answer"`
}

func (e *ElasticSearch) CreateIndex() error {
    client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("http://%s:%d", ElasticSearchService.host, ElasticSearchService.port)), elastic.SetSniff(false))
    if err != nil {
        log.Printf("Error creating Elasticsearch client: %v", err)
		return err
    }

    info, code, err := client.Ping(fmt.Sprintf("http://%s:%d", ElasticSearchService.host, ElasticSearchService.port)).Do(context.Background())
    if err != nil {
        log.Printf("Elasticsearch is not reachable: %v", err)
		return err
    }
    log.Printf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

    exists, err := client.IndexExists("faq").Do(context.Background())
    if err != nil {
        log.Fatalf("Error checking if index exists: %v", err)
		return err
    }

    if !exists {
        // Создаем индекс
        mapping := `{
            "settings": {
                "number_of_shards": 1,
                "number_of_replicas": 0
            },
            "mappings": {
                "properties": {
                    "question": { "type": "text" },
                    "answer": { "type": "text" }
                }
            }
        }`
        _, err := client.CreateIndex("faq").BodyString(mapping).Do(context.Background())
        if err != nil {
            log.Printf("Error creating index: %v", err)
			return err
        }
        log.Println("Index 'faq' created successfully")
        

        faqs := []FAQ{
        {"What is the purpose of your low-code platform?", "Our low-code platform is designed to simplify the creation of applications, enabling users to build solutions quickly without extensive coding knowledge."},
        {"Do I need coding knowledge to use this platform?", "No, the platform is user-friendly and does not require prior coding experience. It’s designed for individuals of all skill levels."},
        {"Can I integrate payment gateways?", "Yes, our platform supports integration with popular payment gateways, allowing you to securely process transactions."},
        {"Is it possible to customize my online store?", "Absolutely! You can customize your online store’s design, layout, and functionality to suit your business needs."},
        {"Does this platform support multiple languages?", "Yes, our platform includes support for multiple languages to help you reach a global audience."},
        {"Can I use my custom domain name?", "Yes, you can connect your custom domain name to your store for a professional and unique web presence."},
        {"How secure is my online store?", "Your online store is protected with industry-standard security features, including SSL encryption and secure payment processing."},
        {"What happens if I encounter technical issues?", "You can contact our support team, which is available to assist you with technical issues and provide solutions."},
        {"Does the platform have SEO tools?", "Yes, our platform includes built-in SEO tools to help improve your store’s visibility on search engines."},
        {"Can I sell digital products with this platform?", "Yes, the platform supports selling digital products, such as eBooks, software, and online courses."},
        {"How much does the platform cost?", "The platform offers flexible pricing plans based on your needs. Contact our sales team for detailed pricing information."},
        {"Are there hidden fees for transactions?", "No, there are no hidden fees for transactions. All costs are transparent and outlined in the pricing plans."},
        {"Can I manage inventory through the platform?", "Yes, the platform includes inventory management tools to help you keep track of stock levels and product availability."},
        {"Is there a mobile version of the store?", "Yes, your online store will be fully responsive and optimized for mobile devices."},
        {"Can I add promotional offers?", "Certainly! You can create and manage promotional offers to attract customers and boost sales."},
        {"Does the platform support social media integration?", "Yes, you can integrate your store with popular social media platforms to enhance your marketing efforts."},
        {"How do I analyze my store’s performance?", "The platform includes analytics tools that provide insights into sales, customer behavior, and website traffic."},
        {"Is there a limit to the number of products I can add?", "No, you can add an unlimited number of products to your store."},
        {"Can I export data from my store?", "Yes, you can export your store data in various formats for backup or migration purposes."},
        {"Does the platform offer email marketing features?", "Yes, you can use the platform’s email marketing tools to communicate with your customers and promote your products."},
        {"Can I set up multiple storefronts?", "Yes, the platform allows you to create and manage multiple storefronts from a single account."},
        {"What payment methods are supported?", "Our platform supports a variety of payment methods, including credit cards, PayPal, and other online payment options."},
        {"Does the platform include shipping management?", "Yes, you can manage shipping settings, rates, and carriers directly through the platform."},
        {"Can I use third-party integrations?", "Absolutely! You can integrate third-party apps and services to extend the functionality of your store."},
        {"What happens if I want to scale my store?", "The platform is designed to grow with your business, offering scalability for increased traffic and sales."},
        {"How long does it take to set up a store?", "You can set up your store within minutes using our intuitive setup wizard."},
        {"Can I upload product images?", "Yes, you can upload high-quality images for your products to showcase them effectively."},
        {"Does the platform support product variations?", "Yes, you can create and manage product variations, such as size, color, and style."},
        {"How do I handle returns and refunds?", "The platform includes tools to manage return and refund processes smoothly and efficiently."},
        {"Can I offer subscription-based products?", "Yes, the platform supports subscription-based products, allowing you to set recurring billing options."},
        {"What type of technical support is available?", "Our dedicated support team provides technical assistance via email, chat, and phone."},
        {"Can I schedule product launches?", "Yes, you can schedule product launches and set specific availability dates."},
        {"Does the platform support real-time chat with customers?", "Yes, you can integrate real-time chat tools to communicate directly with your customers."},
        {"Can I add blog content to my store?", "Yes, you can create and manage blog posts to share updates and engage with your audience."},
        {"How do I handle taxes for my sales?", "The platform provides tax management features to help you calculate and apply sales taxes automatically."},
        {"Can I test my store before launching?", "Yes, you can preview and test your store before it goes live to ensure everything is set up correctly."},
        {"Does the platform provide hosting?", "Yes, the platform includes hosting, so you don’t need to worry about finding a separate hosting provider."},
        {"Can I track abandoned carts?", "Yes, you can track abandoned carts and send follow-up emails to encourage customers to complete their purchases."},
        {"Does the platform support multiple currencies?", "Yes, you can set up your store to support multiple currencies for international customers."},
        {"Is there a way to get feedback from customers?", "Yes, you can collect customer feedback through surveys, reviews, and ratings."},
        {"Can I showcase featured products?", "Yes, you can highlight featured products on your store’s homepage to attract more attention."},
        {"Does the platform support affiliate programs?", "Yes, you can set up affiliate programs to incentivize others to promote your products."},
        {"How do I back up my store data?", "The platform offers tools to back up your store data for added security and peace of mind."},
        {"Can I set up custom shipping zones?", "Yes, you can define custom shipping zones and rates based on customer locations."},
        {"Does the platform offer loyalty programs?", "Yes, you can create loyalty programs to reward customers and encourage repeat purchases."},
        {"How can I optimize my store’s speed?", "The platform includes tools to help optimize your store’s loading times and overall performance."},
        {"What file formats are supported for uploads?", "The platform supports multiple file formats for images, documents, and other media uploads."},
        {"Can I integrate with accounting software?", "Yes, you can connect your store to popular accounting software for easy bookkeeping."},
        {"How do I ensure my store meets legal requirements?", "The platform provides templates and guidelines to help you comply with legal and regulatory requirements."},
        {"Can I automate repetitive tasks?", "Yes, you can automate repetitive tasks such as inventory updates and customer communications using built-in tools."},
    }

        for _, faq := range faqs {
            _, err := client.Index().
                Index("faq").
                BodyJson(faq).
                Do(context.Background())
            if err != nil {
                log.Printf("Failed to insert document: %v", err)
            } else {
                log.Printf("Inserted question: %s", faq.Question)
            }
        }

        log.Print("All FAQ documents added successfully")
    }
    
	return nil
}

func (s *ElasticSearch) GetFaq(keywords *proto.KeywordResponse) ([]FAQ, error) {
    client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("http://%s:%d", s.host, s.port)), elastic.SetSniff(false))
    if err != nil {
        log.Printf("Error creating Elasticsearch client: %v", err)
		return nil, err
    }

    searchResult, err := client.Search().
        Index("faq").
        Query(elastic.NewMultiMatchQuery(strings.Join(keywords.Keywords, " "), "question", "answer")).
        Do(context.Background())
    if err != nil {
        log.Printf("Error executing search query: %v", err)
		return nil, err
    }

    var faqs []FAQ
    for _, hit := range searchResult.Hits.Hits {
        var faq FAQ
        if err := json.Unmarshal(hit.Source, &faq); err != nil {
            log.Printf("Error unmarshalling document: %v", err)
            continue
        }
        faqs = append(faqs, faq)

        if len(faqs) == 5 {
            break
        }
    }

	return faqs, nil
}