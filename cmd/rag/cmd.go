package main

import (
	"fmt"
	"github.com/yonidavidson/gophercon-israel-2024/prompt"
	"github.com/yonidavidson/gophercon-israel-2024/provider"
	"github.com/yonidavidson/gophercon-israel-2024/rag"
	"os"
)

const promptTemplate = `<system>{{.SystemPrompt}}</system>
<user>
Context: 
{{.RAGContext}}

User Query: {{.UserQuery}}</user>`

type promptData struct {
	MaxTokens    float64
	RAGContext   string
	UserQuery    string
	SystemPrompt string
}

func main() {
	apiKey := os.Getenv("PRIVATE_OPENAI_KEY")
	if apiKey == "" {
		fmt.Println("Error: PRIVATE_OPENAI_KEY environment variable not set")
		return
	}
	p := provider.OpenAIProvider{APIKey: apiKey}
	r := rag.New(p)
	es, err := r.Embed(txt, 1000)
	if err != nil {
		fmt.Printf("Error embedding text: %v\n", err)
		return
	}
	fmt.Printf("Number of embeddings: %d\n", len(es))
	userQuery := "What where the conclusions of the research?"
	ragContext, err := r.Search(userQuery, es)
	if err != nil {
		fmt.Printf("Error searching text: %v\n", err)
		return
	}
	m, err := prompt.ParseMessages(promptTemplate, promptData{
		MaxTokens:    1000,
		RAGContext:   string(ragContext),
		UserQuery:    userQuery,
		SystemPrompt: "Answer the following question based only on the provided context:",
	})
	if err != nil {
		fmt.Printf("Error parsing messages: %v\n", err)
		return
	}

	c, err := p.ChatCompletion(m)
	if err != nil {
		fmt.Printf("Error getting chat completion: %v\n", err)
		return
	}
	fmt.Println("\n\n\n\n" + string(c))
}

var txt = `
Title: Using Deep Learning Techniques for Classifications of Radio Signals
Author: Yoni Davidson
Degree: Master of Science in Engineering
Institution: Ariel University Faculty of Electrical Engineering
Date: September 2021

Abstract:
Radio transmission detection and classification without prior signal knowledge is a well-known challenge that has been researched for over 40 years. The evolution of artificial intelligence has opened new avenues for addressing this problem, particularly with the advent of machine learning techniques that can potentially achieve results closer to Shannon’s theoretical limits. The interest in this field has surged again, driven by the need to shift logic to edge devices like IoT and cellphones and to offload logic from hardware to software, enhancing flexibility for various applications, including military communications and radio interference detection.

The objectives of this research were twofold: first, to fine-tune and improve AI methods for classifying radio transmissions, and second, to develop a new generic tool for solving this classification problem. This study integrates mathematical domain knowledge with machine learning techniques, optimizing class grouping to enhance performance. The results demonstrated improved classification accuracy without the need for additional data or changes to the neural network architecture.

Introduction:

1.1 Literature Review:

Analog Modulation:
Analog modulation techniques are essential in transferring low-frequency baseband signals over higher-frequency signals like radio frequencies. Key methods include:

Amplitude Modulation (AM): Modifies the amplitude of the carrier signal to match the modulating signal. Variants include Double Side Band Suppressed Carrier (DSBSC), Single Sideband Suppressed Carrier (SSBSC), and Vestigial Sideband Amplitude Modulation (VSBAM).
Angular Modulation: Involves modulating the frequency or phase of the carrier signal, leading to Frequency Modulation (FM) and Phase Modulation (PM).
Frequency Modulation (FM): Encodes information by varying the instantaneous frequency of the carrier wave.
Phase Modulation (PM): Modulates the phase of the carrier wave to follow the amplitude of the message signal.
Digital Modulation:
Digital Modulation (DM) uses discrete signals to modulate a carrier wave. The three main types are:

Frequency Shift Keying (FSK): Transmits digital information through discrete frequency changes in the carrier signal.
Phase Shift Keying (PSK): Modifies the phase of the carrier signal by varying sine and cosine inputs.
Amplitude Shift Keying (ASK): Represents digital data as variations in the amplitude of the carrier wave.
Modulation Recognition:
Automatic modulation recognition has gained significant attention, particularly in military and academic research. It is a crucial component of communication intelligence (COMINT) systems, which typically include a receiver front-end, a modulation recognizer, and an output stage. The recognition process often begins with signal demodulation, requiring accurate knowledge of the signal's modulation type.

Deep Learning:
Deep learning, a subset of machine learning, involves using neural networks to learn data representations through multiple layers of abstraction. In image recognition, for example, deep learning models might first recognize edges, then patterns, and eventually entire objects like faces. This layered approach, called an artificial neural network (ANN), is also effective in signal processing tasks.

Neural Networks:
Artificial neural networks (ANNs) are computational systems inspired by biological neural networks. They consist of interconnected layers of nodes (neurons) that process input data and generate outputs. ANNs are versatile and have been applied to tasks ranging from image classification to natural language processing.

Convolutional Neural Networks (CNNs):
CNNs are a type of deep learning model specifically designed for image analysis. They use convolutional layers to automatically learn spatial hierarchies of features, making them particularly effective for tasks like image and video recognition, medical image analysis, and more.

Bayes’ Theorem:
Bayes’ theorem is a fundamental concept in probability theory, describing how to update the probability of an event based on new evidence. In the context of machine learning, Bayesian inference, a method derived from Bayes’ theorem, is often used to improve model predictions by incorporating prior knowledge.

XGBoost:
XGBoost is a powerful gradient boosting framework that has gained popularity for its efficiency and performance in machine learning competitions. It is particularly well-suited for problems with a rich feature set and a limited number of samples.

1.2 Objective:
The goal of this research was to enhance the accuracy of a neural network model used for radio signal classification. By integrating domain knowledge into the model, the study aimed to improve classification performance without increasing data size or modifying the network's architecture.

1.3 Organization of the Thesis:
The thesis is organized into five chapters. Chapter 1 introduces the problem and reviews the relevant literature. Chapter 2 outlines the methodology used to improve classification accuracy, including the application of Bayesian statistics. Chapter 3 presents the experimental results, comparing the performance of base and grouped models. Chapter 4 discusses the findings, and Chapter 5 concludes with a summary of results and recommendations for future research.

Methodology:

2.1 Accuracy Improvement by Bayesian Network:
The research builds on a previous study that achieved an 84.7% accuracy rate in classifying 11 different modulations using a basic neural network. The dataset used included 8 digital and 3 analog modulations commonly found in wireless communication systems. The signals were modulated at a rate of approximately 8 samples per symbol with a normalized average transmit power of 0dB. The dataset was analyzed in the time domain, and visual similarities between different modulations were observed.

The classification errors were analyzed using a confusion matrix, which revealed that certain modulations were often confused with one another. Based on these findings, the modulations were grouped into categories to improve classification accuracy. For example, QAM16 and QAM64 were grouped together, as were AM-DSB, AM-SSB, and WBFM.

The grouped modulations were then classified using machine learning techniques based on statistical features. The Python library tsfresh was used to automatically calculate a large number of time-series characteristics, and the XGBoost algorithm was employed to classify modulations within each group. The results indicated that this approach significantly improved classification accuracy, particularly for modulations that were difficult to distinguish in the original model.

2.2 Differentiating Between Modulations in the Same Groups:
To further refine the classification process, statistical features were extracted from the grouped modulations using tsfresh. XGBoost was then applied to classify modulations within each group. The features that had the most significant impact on classification accuracy were identified using the eli5 library, which provides an explanation of the model's decisions.

For example, the "value__absolute_sum_of_changes" feature was found to be particularly important in distinguishing between CPFSK and GFSK modulations. By analyzing these features, the researchers were able to gain insights into the underlying characteristics of each modulation and improve the classification model accordingly.

Results:

3.1 Comparing Base and Grouped Models:

3.1.1 Training Properties Comparison:
The research compared the performance of the Classification Base Model (CBM) and the Classification Grouped Model (CGM). The CGM outperformed the CBM in several key areas:

Training Efficiency: The CGM required fewer resources for training, with each epoch taking 25% less time than the CBM.
Initial Error Reduction: The CGM showed a smaller error rate from the beginning of training, indicating that the grouping approach was effective.
Sensitivity: The CGM achieved a 70% detection rate at -7dB, whereas the CBM required 0dB to reach the same accuracy level.
3.1.2 Model Results for Signal Lower than Detection Level:
Both models struggled to detect signals with a Signal-to-Noise Ratio (SNR) lower than -16dB. The system tended to classify noise as AM-SSB, which was expected given the characteristics of the Gaussian noise generated for the experiment.

3.1.3 Model Results for SNR Close to Detection:
At SNR values of -8dB and higher, the grouped model began to detect modulations more accurately. The CGM showed improved detection rates compared to the CBM, particularly for QAM modulations.

3.1.4 Increase with Detection Rate in Correlation to SNR -6dB to -2dB:
Both models showed a linear increase in detection rate as the SNR improved. However, the CGM consistently outperformed the CBM across this SNR range.

3.1.5 Increase with Detection Rate in Correlation to SNR 0dB to 4dB:
As the SNR increased from 0dB to 4dB, the CGM reached near-maximum detection rates, while the CBM continued to show linear improvement.

3.1.6 Increase with Detection Rate in Correlation to SNR 6dB to 10dB:
At higher SNR values, the CGM maintained its superior performance, achieving maximum detection rates faster than the CBM.

3.1.7 Increase with Detection Rate in Correlation to SNR 12dB to 16dB:
Both models reached their steady-state detection rates in this SNR range, with the CGM consistently achieving better results.

3.1.8 Best Sample Comparison SNR 18dB:
At 18dB, both models reached their full detection rates. However, the CGM was more effective at distinguishing between modulations that were challenging for the CBM, such as QAM16, QAM64, and 8PSK.

3.1.9 Final Receiver Diagram:
The final results showed that the CGM achieved a 5dB improvement in receiver sensitivity compared to the CBM. The grouped approach allowed the same neural network model to achieve higher accuracy without additional training data or changes to the network architecture.

3.2 Comparing in a Subgroup:

3.2.1 QAM16 vs. QAM64:
The CGM achieved 93.48% accuracy in distinguishing between QAM16 and QAM64 modulations. Key features included the "value__ar_coefficient" and "value__change_quantiles" metrics.

3.2.2 8PSK vs. BPSK vs. QPSK:
The CGM achieved 82.63% accuracy for this subgroup, with important features including the "value__minimum" and "value__abs_energy."

3.2.3 AM-DSB vs. AM-SSB vs. WBFM:
The CGM achieved 82.22% accuracy, with the "value__standard_deviation" and "value__absolute_sum_of_changes" being the most influential features.

3.2.4 CPFSK vs. GFSK:
The CGM achieved a perfect 100% accuracy rate in distinguishing between CPFSK and GFSK, with the "value__absolute_sum_of_changes" feature being the most significant.

Discussion:

4.1 Review of Main Findings:
The research demonstrated that grouping modulations based on confusion matrix results and domain knowledge significantly improves classification performance. The CGM outperformed the CBM in both accuracy and sensitivity, making it a more effective tool for practical applications in radio signal classification. By leveraging prior knowledge and optimizing class grouping, the researchers achieved a 10% improvement in detection success rate.

This approach also reduced training time, making the CGM more efficient than the CBM. The methodology used in this study can be applied to other domains where signal classification is required, particularly in scenarios with limited data or computational resources.

Conclusion:

5.1 Summary of Results:
The study successfully increased the detection rate by 10% without increasing the dataset size or altering the neural network architecture. The research highlights the importance of combining deep learning with domain knowledge to solve complex engineering problems. The methodology developed in this study provides a framework for improving classification performance in various signal processing applications.

5.2 Recommendations for Further Studies:
Future research should explore the application of this methodology to edge devices, where limited memory and energy constraints require efficient solutions. The approach could also be adapted for hybrid systems that combine deep learning with statistical methods to optimize performance on resource-constrained devices.

Bibliography:
The thesis references a range of studies on deep learning, modulation techniques, and signal processing, providing a comprehensive foundation for the research.
`
