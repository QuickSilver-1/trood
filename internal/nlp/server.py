import grpc
from concurrent import futures
from proto import server_pb2 as pb2
from proto import server_pb2_grpc as pb2_grpc
import spacy

nlp = spacy.load("model")

class KeywordExtractorService(pb2_grpc.KeywordExtractorServicer):
    def ExtractKeywords(self, request, context):
        question = request.question
        doc = nlp(question)
        keywords = [token.text for token in doc if token.is_alpha and not token.is_stop]
        
        return pb2.KeywordResponse(keywords=keywords)

def start():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb2_grpc.add_KeywordExtractorServicer_to_server(KeywordExtractorService(), server)
    server.add_insecure_port("[::]:50051")
    
    print("gRPC server is running on port 50051...")
    server.start()
    server.wait_for_termination()
