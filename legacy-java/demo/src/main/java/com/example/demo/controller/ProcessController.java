package com.example.demo.controller;

import com.example.demo.model.ProcessRequest;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class ProcessController {
    
    @PostMapping("/process")
    public ProcessRequest process(@RequestBody ProcessRequest request) {
        System.out.println("Received request: " + request.getMessage());
        return request;
    }
}